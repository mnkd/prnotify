package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/m-nakada/slackposter"
)

const (
	ExitCodeOK        int = iota // 0
	ExitCodeError                // 1
	ExitCodeFileError            // 2
)

type App struct {
	Config    Config
	Slack     slackposter.Slack
	GitHubAPI GitHubAPI
}

var app = App{}

func (app App) Run() int {
	pulls, err := app.GitHubAPI.GetPulls()
	if err != nil {
		return ExitCodeError
	}

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

func init() {
	var slackConfigIndex int
	var owner, repo, accessToken string
	flag.StringVar(&owner, "owner", "", "owner")
	flag.StringVar(&repo, "repo", "", "GitHub repository")
	flag.StringVar(&accessToken, "token", "", "GitHub API Access Token")
	flag.IntVar(&slackConfigIndex, "sci", 0, "Slack Config Index (default: 0)")
	flag.Parse()

	if len(owner) == 0 || len(repo) == 0 || len(accessToken) == 0 {
		flag.Usage()
		os.Exit(ExitCodeError)
	}

	// Prepare config
	config, err := NewConfig()
	if err != nil {
        os.Exit(ExitCodeError)
	}

	// Prepare GitHubAPI
	gh := NewGitHubAPI()
	gh.AccessToken = accessToken
	gh.Owner = owner
	gh.Repo = repo

	// Steup app
	app.Config = config
	app.GitHubAPI = gh
	app.Slack = slackposter.NewSlack(config.SlackChannels[slackConfigIndex])

	// Setup users map
	app.GitHubAPI.UsersMap, err = NewUsers("users.json")
	if err != nil {
        os.Exit(ExitCodeError)
	}
}

func main() {
	os.Exit(app.Run())
}
