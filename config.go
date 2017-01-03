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
	Config_InvalidJsonError               = errors.New("Config_InvalidJsonError")
	Config_InvalidSlackWebhooksIndexError = errors.New("Config_InvalidSlackWebhooksIndexError")
)

type Config struct {
	GitHub struct {
		AccessToken string `json:"access_token"`
		Owner       string `json:"owner"`
		Repo        string `json:"repo"`
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
		return Config_InvalidJsonError
	}
	if len(config.GitHub.Owner) == 0 {
		fmt.Fprintln(os.Stderr, "Invalid config.json. You should set a github owner.")
		return Config_InvalidJsonError
	}
	if len(config.GitHub.Repo) == 0 {
		fmt.Fprintln(os.Stderr, "Invalid config.json. You should set a github repo.")
		return Config_InvalidJsonError
	}

	if len(config.SlackWebhooks) < config.SlackWebhooksIndex+1 {
		fmt.Fprintln(os.Stderr, "Invalid slack webhooks index:", config.SlackWebhooksIndex)
		return Config_InvalidSlackWebhooksIndexError
	}

	return nil
}

func NewConfig(path string, slackWebhooksIndex int, dryRun bool) (Config, error) {
	var config Config
	config.SlackWebhooksIndex = slackWebhooksIndex
	config.DryRun = dryRun

	usr, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Could not get current user:", err)
		return config, err
	}

	if len(path) == 0 {
		path = filepath.Join(usr.HomeDir, "/.config/prnotify/config.json")
	} else {
		p, err := filepath.Abs(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "[Error] Could not return absolute representation of path:", err, path)
			return config, err
		}
		path = p
	}

	str, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Could not read config.json:", err)
		return config, err
	}

	if err := json.Unmarshal(str, &config); err != nil {
		fmt.Fprintln(os.Stderr, "[Error] JSON Unmarshal:", err)
		return config, err
	}

	if err = config.validate(); err != nil {
		return config, err
	}

	return config, nil
}
