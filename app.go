package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/m-nakada/slackposter"
)

type App struct {
	Config       Config
	Slack        slackposter.Slack
	GitHubAPI    GitHubAPI
	UsersManager UsersManager
}

func (app App) ActivePullRequests() ([]PullRequest, error) {
	pulls, err := app.GitHubAPI.GetPulls()
	if err != nil {
		return pulls, err
	}

	var activePulls []PullRequest
	for _, pull := range pulls {
		if strings.Contains(strings.ToUpper(pull.Title), "WIP") {
			continue
		}
		activePulls = append(activePulls, pull)
	}
	return activePulls, nil
}

func (app App) Run() int {

	// Build Payload
	builder := NewMessageBuilder(app.GitHubAPI, app.UsersManager)

	var payload slackposter.Payload
	payload.Channel = app.Slack.Channel
	payload.Username = app.Slack.Username
	payload.IconEmoji = app.Slack.IconEmoji
	payload.LinkNames = true
	payload.Mrkdwn = true

	pulls, err := app.ActivePullRequests()
	if err != nil {
		return ExitCodeError
	}

	payload.Text = builder.BudildSummary(len(pulls))

	var attachments []slackposter.Attachment
	for _, pull := range pulls {
		fmt.Fprintf(os.Stdout, "#%d\n", pull.Number)
		comments, err := app.GitHubAPI.GetCommentsWithPullRequest(pull)
		if err != nil {
			return ExitCodeError
		}

		attachment := builder.BuildAttachment(pull, comments)
		attachments = append(attachments, attachment)
	}
	payload.Attachments = attachments

	err = app.Slack.PostPayload(payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Could not send a payload to slack:", err)
		return ExitCodeError
	}

	return ExitCodeOK
}

func NewApp(config Config) (App, error) {
	var app = App{}
	var err error
	app.Config = config

	app.GitHubAPI = NewGitHubAPI(config)
	app.UsersManager, err = NewUsersManager("users.json")
	app.Slack = slackposter.NewSlack(config.SlackWebhooks[config.SlackWebhooksIndex])
	// app.Slack.DryRun = true

	return app, err
}
