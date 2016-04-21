package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/andreaskoch/ga-spam-control/api"
	"github.com/andreaskoch/ga-spam-control/cli/templates"
	"github.com/andreaskoch/ga-spam-control/spamcontrol"
	"github.com/mitchellh/go-homedir"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	handleCommandlineArguments(os.Args[1:])
}

// handleCommandlineArguments parses the given arguments
// and performs the selected action.
func handleCommandlineArguments(args []string) {
	app := kingpin.New("ga-spam-control", "Command-line utility for blocking referer spam from your Google Analytics accounts")
	app.Version("0.0.1")

	status := app.Command("status", "Display the current spam control status of your accounts")
	update := app.Command("update", "Update your spam control settings")
	remove := app.Command("remove", "Remove spam control from your accounts")

	switch kingpin.MustParse(app.Parse(args)) {

	// Display status
	case status.FullCommand():
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		statusError := cli.Status()
		if statusError != nil {
			app.Fatalf("%s", statusError.Error())
		}

		os.Exit(0)

	// Update filters
	case update.FullCommand():
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		updateError := cli.Update()
		if updateError != nil {
			app.Fatalf("%s", updateError.Error())
		}

		os.Exit(0)

	// Remove spam control
	case remove.FullCommand():
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		removeError := cli.Remove()
		if removeError != nil {
			app.Fatalf("%s", removeError.Error())
		}

		os.Exit(0)

	}

	os.Exit(0)
}

// newCLI creates a new spam control instance.
func newCLI() (*cli, error) {

	// create a token store
	homeDirPath, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("Cannot determine the current users home direcotry. Error: %s", err)
	}

	tokenStoreFilePath := filepath.Join(homeDirPath, ".analytics")
	tokenStore := newTokenStore(tokenStoreFilePath)

	// create a new analytis API instance
	googleAnalyticsClientID := "367149948041-v2up5mcsv4a415gm9rmmmli5lifucddr.apps.googleusercontent.com"
	googleAnalyticsClientSecret := "P9GswEyyNGUuewVebN2k6EjH"
	analyticsAPI, apiError := api.New(tokenStore, googleAnalyticsClientID, googleAnalyticsClientSecret)
	if apiError != nil {
		return nil, apiError
	}

	// create a spam control instance
	spamControl := spamcontrol.New(analyticsAPI)

	return &cli{
		spamControl: spamControl,
	}, nil

}

// cli contains functions for managing
// spam control filters of Google Analytics accounts.
type cli struct {
	spamControl spamcontrol.SpamController
}

// Update the spam control filters.
func (cli *cli) Update() error {
	return cli.spamControl.Update()
}

// Remove all spam control filters.
func (cli *cli) Remove() error {
	return cli.spamControl.Remove()
}

// Status displays the spam control status.
func (cli *cli) Status() error {
	statusViewModel, err := cli.spamControl.Status()
	if err != nil {
		return err
	}

	statusTemplate, parseError := template.New("Status").Parse(templates.Status)
	if parseError != nil {
		return parseError
	}

	renderError := statusTemplate.Execute(os.Stdout, statusViewModel)
	if renderError != nil {
		return renderError
	}

	return nil
}
