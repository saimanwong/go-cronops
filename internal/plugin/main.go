package plugin

import (
	"github.com/saimanwong/go-cronops/internal/plugin/example"
	"github.com/saimanwong/go-cronops/internal/plugin/jiraslack"
)

func NewExample() *example.Example {
	return &example.Example{}
}

func NewJiraSlack() *jiraslack.JiraSlack {
	return &jiraslack.JiraSlack{}
}
