package jiraslack

import (
	"testing"

	"github.com/mitchellh/mapstructure"
)

const (
	logprefix    = "SlackJira"
	jiraUrl      = "localhost:8080"
	jiraUsername = "test"
	jiraPassword = "test"
)

func TestParam(t *testing.T) {
	param := map[string]interface{}{
		"logprefix":    logprefix,
		"jiraUrl":      jiraUrl,
		"jiraUsername": jiraUsername,
		"jiraPassword": jiraPassword,
	}
	var js JiraSlack
	err := mapstructure.Decode(param, &js)
	if err != nil {
		t.Errorf("could not decode param")
	}

	if js.LogPrefix != logprefix {
		t.Errorf("LogPrefix and logprefix does not match")
	}

	if js.JiraUrl != jiraUrl {
		t.Errorf("JiraUrl and jira_url does not match")
	}

	if js.JiraUsername != jiraUsername {
		t.Errorf("JiraUser and jira_username does not match")
	}

	if js.JiraPassword != jiraPassword {
		t.Errorf("JiraPass and jira_password does not match")
	}
}
