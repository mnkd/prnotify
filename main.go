package main

import (
	"flag"
	"os"
)

const (
	ExitCodeOK int = iota
	ExitCodeError
	ExitCodeFileError
)

var app App

func init() {
	var configPath string
	var index int
	var dryRun bool
	flag.StringVar(&configPath, "c", "", "/path/to/config.json. (default: $HOME/.config/prnotify/config.json)")
	flag.IntVar(&index, "swi", 0, "Slack Webhooks Index (default: 0)")
	flag.BoolVar(&dryRun, "d", false, "A dry run will not send any message to Slack. (defualt: false)")
	flag.Parse()

	// Prepare config
	config, err := NewConfig(configPath, index, dryRun)
	if err != nil {
		os.Exit(ExitCodeError)
	}

	// Prepare app
	app, err = NewApp(config)
	if err != nil {
		os.Exit(ExitCodeError)
	}
}

func main() {
	os.Exit(app.Run())
}
