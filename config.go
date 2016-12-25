package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/m-nakada/slackposter"
)

type Config struct {
	SlackChannels []slackposter.Config `json:"channels"`
}

func NewConfig() (Config, error) {
	var config Config

	usr, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Could not get current user:", err)
		return config, err
	}

	path := filepath.Join(usr.HomeDir, "/.config/slackposter/config.json")
	str, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Could not read config.json:", err)
		return config, err
	}

	if err := json.Unmarshal(str, &config); err != nil {
		fmt.Fprintln(os.Stderr, "[Error] JSON Unmarshal:", err)
		return config, err
	}

	return config, nil
}
