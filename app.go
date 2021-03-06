package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mnkd/dayoff"
	"github.com/mnkd/holidayJP"
	"github.com/mnkd/slackposter"
)

type App struct {
	Config       Config
	Slack        slackposter.SlackPoster
	GitHubAPI    GitHubAPI
	UsersManager UsersManager
}

// Attachments index
type AttachmentType int

const (
	MERGE AttachmentType = iota
	REVIEW
	CHECK
	ASSIGN_REVIEWER
)

func isHoliday() bool {
	return holidayJP.IsHoliday(time.Now())
}

func isDayOff() bool {
	return dayoff.IsDayOff(time.Now())
}

func (app App) ActivePullRequests() ([]PullRequest, error) {
	var activePulls []PullRequest

	pulls, err := app.GitHubAPI.GetPulls()
	if err != nil {
		return activePulls, err
	}

	for _, pull := range pulls {
		if strings.Contains(strings.ToUpper(pull.Title), "WIP") {
			continue
		}
		activePulls = append(activePulls, pull)
	}
	return activePulls, nil
}

func (app App) Run() int {
	if app.Config.App.UseHolidayJP && isHoliday() {
		fmt.Fprintln(os.Stderr, "It's a holiday.")
		return ExitCodeError
	}

	if app.Config.App.UseDayOff && isDayOff() {
		fmt.Fprintln(os.Stderr, "It's a day off.")
		return ExitCodeError
	}

	// Build Payload
	builder := NewMessageBuilder(app.GitHubAPI, app.UsersManager, app.Config)

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

	// Prepare summary
	summary := builder.BudildSummary(len(pulls))
	payload.Text = summary
	fmt.Fprintf(os.Stdout, "\n%v", summary)

	// Prepare fields for each attachment
	fieldsMap := make(map[AttachmentType][]slackposter.Field)
	for i, pull := range pulls {
		fmt.Fprintf(os.Stdout, "%-2d #%d\n", i+1, pull.Number)

		reviews, err := app.GitHubAPI.GetReviewsWithPullRequest(pull)
		if err != nil {
			return ExitCodeError
		}

		field, attachmentType := builder.BuildField(pull, reviews)
		fieldsMap[attachmentType] = append(fieldsMap[attachmentType], field)
	}

	// Prepare attachments
	var attachments []slackposter.Attachment
	for i := MERGE; i < ASSIGN_REVIEWER+1; i++ {
		if len(fieldsMap[i]) == 0 {
			continue
		}

		attachment := builder.BuildAttachmentWithType(i)
		attachment.Fields = fieldsMap[i]

		attachments = append(attachments, attachment)
	}
	payload.Attachments = attachments

	// Post payload
	err = app.Slack.PostPayload(payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "App: <error> send a payload to slack:", err)
		return ExitCodeError
	}

	return ExitCodeOK
}

func NewApp(config Config) (App, error) {
	var app = App{}
	var err error
	app.Config = config

	app.GitHubAPI = NewGitHubAPI(config)
	app.UsersManager, err = NewUsersManager()
	app.Slack = slackposter.NewSlackPoster(config.SlackWebHook)
	app.Slack.DryRun = config.DryRun

	return app, err
}
