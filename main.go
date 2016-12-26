package main

import (
	"flag"
	"os"
)

const (
	ExitCodeOK        int = iota // 0
	ExitCodeError                // 1
	ExitCodeFileError            // 2
)

var app App

func init() {
	var index int
	var rich bool
	var configPath string
	flag.IntVar(&index, "swi", 0, "Slack Webhooks Index (default: 0)")
	flag.BoolVar(&rich, "rich", false, "Use rich format (default: false)")
	flag.StringVar(&configPath, "c", "", "/path/to/config.json. (default: $HOME/.config/prnotify/config.json)")
	flag.Parse()

	// Prepare config
	config, err := NewConfig(configPath, index, rich)
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
