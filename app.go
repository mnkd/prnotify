package main

import (
	"fmt"
	"os"

	"github.com/m-nakada/slackposter"
)

type App struct {
	Config    Config
	Slack     slackposter.Slack
	GitHubAPI GitHubAPI
}

func (app App) Run() int {
	pulls, err := app.GitHubAPI.GetPulls()
	if err != nil {
		return ExitCodeError
	}

	if app.Config.RichFormat {
		payload := app.GitHubAPI.SlackPayload(pulls)
		payload.Channel = app.Slack.Channel
		payload.Username = app.Slack.Username
		payload.IconEmoji = app.Slack.IconEmoji
		payload.LinkNames = true

		err = app.Slack.PostPayload(payload)
		if err != nil {
			fmt.Fprintln(os.Stderr, "[Error] Could not send a payload to slack:", err)
			return ExitCodeError
		}

		return ExitCodeOK
	}

	// Plain text
	message := app.GitHubAPI.SlackMessage(pulls)
	if len(message) == 0 {
		fmt.Fprintln(os.Stdout, "No message.")
		return ExitCodeOK
	}

	err = app.Slack.PostMessage(message)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Could not send a message to slack:", err)
		return ExitCodeError
	}

	return ExitCodeOK
}

func NewApp(config Config) (App, error) {
	var app = App{}
	var err error
	app.Config = config

	app.GitHubAPI = NewGitHubAPI(config)
	app.GitHubAPI.UsersMap, err = NewUsers("users.json")
	app.Slack = slackposter.NewSlack(config.SlackWebhooks[config.SlackWebhooksIndex])

	return app, err
}
