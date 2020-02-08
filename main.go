package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	ExitCodeOK int = iota
	ExitCodeError
	ExitCodeFileError
)

var (
	Version  string
	Revision string
)

var app App

func init() {
	var configPath string
	var dryRun, version bool
	flag.StringVar(&configPath, "c", "", "/path/to/config.json. (default: $HOME/.config/prnotify/config.json)")
	flag.BoolVar(&dryRun, "d", false, "A dry run will not send any message to Slack. (defualt: false)")
	flag.BoolVar(&version, "v", false, "Print version.")
	flag.Parse()

	if version {
		fmt.Fprintln(os.Stdout, "Version:", Version)
		fmt.Fprintln(os.Stdout, "Revision:", Revision)
		os.Exit(ExitCodeOK)
	}

	// Prepare config
	config, err := NewConfig(configPath, dryRun)
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
