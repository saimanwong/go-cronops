package jiraslack

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"

	l "github.com/saimanwong/go-cronops/internal/logger"

	"github.com/andygrunwald/go-jira"
	"github.com/slack-go/slack"
)

var (
	logger = l.New()
	log    = l.NewLogEntry()
)

type JiraSlack struct {
	LogPrefix           string
	JiraUrl             string
	JiraUsername        string
	JiraPassword        string
	JiraJQL             string
	JiraIgnoreIssues    []string
	JiraIgnoreSummaries []string

	SlackToken           string
	SlackChannelID       string
	SlackFirstMsgTitle   string
	SlackFirstMsgText    string
	SlackAttachmentColor string
}

func (js *JiraSlack) Run() {
	log = logger.WithField("prefix", js.LogPrefix)
	err := js.validate()
	if err != nil {
		log.Errorf("validate input failed, %s", err)
		return
	}

	// Decode (decrypt in the future)
	jp, _ := b64.StdEncoding.DecodeString(js.JiraPassword)
	st, _ := b64.StdEncoding.DecodeString(js.SlackToken)
	jiraPass := strings.Trim(string(jp), "\n")
	slackToken := strings.Trim(string(st), "\n")

	// Create Jira Client
	tp := jira.BasicAuthTransport{
		Username: js.JiraUsername,
		Password: jiraPass,
	}

	jClient, err := jira.NewClient(tp.Client(), js.JiraUrl)
	if err != nil {
		log.Errorf("failed to create jira client, %s", err)
		return
	}

	// Ignore issues before querying
	var jqlBldr strings.Builder
	jqlBldr.WriteString(js.JiraJQL)
	for _, i := range js.JiraIgnoreIssues {
		jqlBldr.WriteString(fmt.Sprintf(" and issueKey!=%s", i))
		log.Infof("ingored issueKey %s", i)
	}

	for _, i := range js.JiraIgnoreSummaries {
		jqlBldr.WriteString(fmt.Sprintf(" and summary!~'%s'", i))
		log.Infof("ignored summary %s", i)
	}

	// Get issues based on JQL
	issues, _, err := jClient.Issue.Search(jqlBldr.String(), nil)
	if err != nil {
		log.Error("failed to get jira issues, ", err)
		return
	}

	// Get slack client
	sClient := slack.New(slackToken)

	// Replace %total_jira_issues% in text
	text := js.SlackFirstMsgText
	re := regexp.MustCompile(`%total_jira_issues%`)
	if re.MatchString(text) {
		text = re.ReplaceAllString(text, fmt.Sprintf("%d", len(issues)))
	}

	// Main slack thread
	ts, err := postMsg(sClient, msg{
		channel:   js.SlackChannelID,
		ts:        nil,
		title:     js.SlackFirstMsgTitle,
		titleLink: "",
		text:      text,
		color:     js.SlackAttachmentColor,
	})
	if err != nil {
		log.Errorf("failed to post message to slack, %T", err)
		return
	}
	log.Infof("posted the first message to slack channel %s", js.SlackChannelID)

	// Posts issues into thread
	for _, i := range issues {
		description := "No description"
		if len(i.Fields.Description) > 0 {
			description = i.Fields.Description
		}
		_, err := postMsg(sClient, msg{
			channel:   js.SlackChannelID,
			ts:        &ts,
			title:     fmt.Sprintf("%s - %s", i.Key, i.Fields.Summary),
			titleLink: fmt.Sprintf("%s%s", js.JiraUrl+"/browse/", i.Key),
			text:      description,
			color:     js.SlackAttachmentColor,
		})
		if err != nil {
			log.Errorf("failed to post message to slack thread, %T", err)
			return
		}
		log.Infof("posted issue %s to slack thread", i.Key)
	}
}

// Function to validate input
func (js *JiraSlack) validate() error {
	reUrl := regexp.MustCompile(`^https?://`)
	if !reUrl.MatchString(js.JiraUrl) {
		return errors.New("jira URL must begin with either http:// or https://")
	}
	strings.TrimRight(js.JiraUrl, "/ ")
	return nil
}

// Wrapper struct for postMsg function.
type msg struct {
	channel   string
	ts        *string
	title     string
	titleLink string
	text      string
	color     string
}

// Wrapper function to post message to slack.
func postMsg(c *slack.Client, m msg) (string, error) {
	var msgOption []slack.MsgOption
	if m.ts != nil {
		msgOption = append(msgOption, slack.MsgOptionTS(*m.ts))
	}
	msgOption = append(msgOption, slack.MsgOptionAttachments(
		slack.Attachment{
			Title:      m.title,
			TitleLink:  m.titleLink,
			Text:       m.text,
			MarkdownIn: []string{"text"},
			Color:      m.color,
		}),
	)

	_, ts, err := c.PostMessage(
		m.channel,
		msgOption...,
	)
	if err != nil {
		return "", err
	}
	return ts, nil
}
