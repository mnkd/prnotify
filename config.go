package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

var (
	ErrInvalidJson               = errors.New("ErrInvalidJson")
	ErrInvalidSlackWebhooksIndex = errors.New("ErrInvalidSlackWebhooksIndex")
)

type Config struct {
	App struct {
		UseHolidayJP bool `json:"use_holiday_jp"`
		UseDayOff    bool `json:"use_dayoff"`
	} `json:"app"`
	GitHub struct {
		AccessToken     string `json:"access_token"`
		Owner           string `json:"owner"`
		Repo            string `json:"repo"`
		MinimumApproved int    `json:"minimum_approved"`
		Comment         struct {
			PerPage int `json:"per_page"`
		} `json:"comment"`
	} `json:"github"`
	SlackWebhooks []struct {
		Channel    string `json:"channel"`
		IconEmoji  string `json:"icon_emoji"`
		Username   string `json:"username"`
		WebhookUrl string `json:"webhook_url"`
	} `json:"slack_webhooks"`

	DryRun             bool
	SlackWebhooksIndex int
}

func (config *Config) validate() error {
	// Validate
	if len(config.GitHub.AccessToken) == 0 {
		fmt.Fprintln(os.Stderr, "Invalid config.json. You should set a github access_token.")
		return ErrInvalidJson
	}
	if len(config.GitHub.Owner) == 0 {
		fmt.Fprintln(os.Stderr, "Invalid config.json. You should set a github owner.")
		return ErrInvalidJson
	}
	if len(config.GitHub.Repo) == 0 {
		fmt.Fprintln(os.Stderr, "Invalid config.json. You should set a github repo.")
		return ErrInvalidJson
	}
	if config.GitHub.Comment.PerPage < 30 {
		fmt.Fprintln(os.Stderr, "Invalid value: github.comment.per_page. (min 30)")
		config.GitHub.Comment.PerPage = 30
	}
	if len(config.SlackWebhooks) < config.SlackWebhooksIndex+1 {
		fmt.Fprintln(os.Stderr, "Invalid slack webhooks index:", config.SlackWebhooksIndex)
		return ErrInvalidSlackWebhooksIndex
	}

	return nil
}

func NewConfig(path string, slackWebhooksIndex int, dryRun bool) (Config, error) {
	var config Config
	config.SlackWebhooksIndex = slackWebhooksIndex
	config.DryRun = dryRun

	usr, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Config: <error> get current user:", err)
		return config, err
	}

	if len(path) == 0 {
		path = filepath.Join(usr.HomeDir, "/.config/prnotify/config.json")
	} else {
		p, err := filepath.Abs(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Config: <error> get absolute representation of path:", err, path)
			return config, err
		}
		path = p
	}

	str, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Config: <error> read config.json:", err)
		return config, err
	}

	if err := json.Unmarshal(str, &config); err != nil {
		fmt.Fprintln(os.Stderr, "Config: <error> json unmarshal: config.json:", err)
		return config, err
	}

	if err = config.validate(); err != nil {
		return config, err
	}

	return config, nil
}
