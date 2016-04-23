package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/andreaskoch/ga-spam-control/api"
	"github.com/andreaskoch/ga-spam-control/cli/templates"
	"github.com/andreaskoch/ga-spam-control/cli/token"
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
	app := kingpin.New("ga-spam-control", "Command-line utility for blocking referrer spam from your Google Analytics accounts")
	app.Version("0.1.0")

	status := app.Command("status", "Display the current spam control status of your accounts")
	statusQuiet := status.Flag("quiet", "Display status in a parsable format").Short('q').Bool()
	stutusAccountID := status.Arg("accountID", "Google Analytics account ID").String()

	update := app.Command("update", "Update your spam control settings")
	updateAccountID := update.Arg("accountID", "Google Analytics account ID").Required().String()

	remove := app.Command("remove", "Remove spam control from your accounts")
	removeAccountID := remove.Arg("accountID", "Google Analytics account ID").Required().String()

	switch kingpin.MustParse(app.Parse(args)) {

	// Display status
	case status.FullCommand():
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		statusError := cli.Status(*stutusAccountID, *statusQuiet)
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

		updateError := cli.Update(*updateAccountID)
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

		removeError := cli.Remove(*removeAccountID)
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

	tokenStoreFilePath := filepath.Join(homeDirPath, ".ga-spam-control")
	tokenStore := token.NewTokenStore(tokenStoreFilePath)

	// create a new analytis API instance
	googleAnalyticsClientID := "821429244906-8aki1tiaov6g2o7lr7elp41435adk9ge.apps.googleusercontent.com"
	googleAnalyticsClientSecret := "_WxLj0SpQ8HxqmOEyYDUTFzW"
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

// Update the spam control filters for the account with the given accountID.
func (cli *cli) Update(accountID string) error {
	return cli.spamControl.Update(accountID)
}

// Remove all spam control filters for the account with the given accountID.
func (cli *cli) Remove(accountID string) error {
	return cli.spamControl.Remove(accountID)
}

// Status displays the spam control status.
func (cli *cli) Status(accountID string, quiet bool) error {

	if accountID == "" {
		return cli.gobalStatus(quiet)
	}

	return cli.accountStatus(accountID)
}

// gobalStatus displays the spam control status of all accessible accounts.
func (cli *cli) gobalStatus(quiet bool) error {
	statusViewModel, err := cli.spamControl.GlobalStatus()
	if err != nil {
		return err
	}

	// select the display template
	templateText := templates.PrettyStatus
	if quiet {
		// the "quiet" template contains no surplus texts
		// and should be easier to parse by tools like awk
		templateText = templates.QuietStatus
	}

	// parse the template
	statusTemplate, parseError := template.New("Status").Parse(templateText)
	if parseError != nil {
		return parseError
	}

	renderError := statusTemplate.Execute(os.Stdout, statusViewModel)
	if renderError != nil {
		return renderError
	}

	return nil
}

// accountStatus displays the spam control status for account with
// the given accountID.
func (cli *cli) accountStatus(accountID string) error {
	accountStatus, err := cli.spamControl.AccountStatus(accountID)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s\n", accountStatus)

	return nil
}
