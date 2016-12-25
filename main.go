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
	flag.IntVar(&index, "swi", 0, "Slack Webhooks Index (default: 0)")
	flag.Parse()

	// Prepare config
	config, err := NewConfig()
	if err != nil {
		os.Exit(ExitCodeError)
	}
	config.SlackWebhooksIndex = index

	// Prepare app
	app, err = NewApp(config)
	if err != nil {
		os.Exit(ExitCodeError)
	}
}

func main() {
	os.Exit(app.Run())
}
