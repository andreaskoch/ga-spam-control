package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/andreaskoch/ga-spam-control/api"
	"github.com/andreaskoch/ga-spam-control/cli/credentials"
	"github.com/mitchellh/go-homedir"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	cli(os.Args[1:])
}

func cli(args []string) {
	app := kingpin.New("ga-spam-control", "Command-line utility for blocking referer spam from your Google Analytics accounts")
	app.Version("0.0.1")

	status := app.Command("status", "Display the current spam control status of your accounts")
	update := app.Command("update", "Update your spam control settings")
	remove := app.Command("remove", "Remove spam control from your accounts")

	switch kingpin.MustParse(app.Parse(args)) {

	// Display status
	case status.FullCommand():
		cli, err := newSpamControl()
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
		cli, err := newSpamControl()
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
		cli, err := newSpamControl()
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

// newSpamControl creates a new spam control instance.
func newSpamControl() (*spamControl, error) {

	// create a token store
	homeDirPath, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("Cannot determine the current users home direcotry. Error: %s", err)
	}

	tokenStoreFilePath := filepath.Join(homeDirPath, ".analytics")
	tokenStore := credentials.NewTokenStore(tokenStoreFilePath)

	// create a new analytis API instance
	googleAnalyticsClientID := "821429244906-8aki1tiaov6g2o7lr7elp41435adk9ge.apps.googleusercontent.com"
	googleAnalyticsClientSecret := "_WxLj0SpQ8HxqmOEyYDUTFzW"
	analyticsAPI, apiError := api.New(tokenStore, googleAnalyticsClientID, googleAnalyticsClientSecret)
	if apiError != nil {
		return nil, apiError
	}

	return &spamControl{
		analyticsAPI: analyticsAPI,
	}, nil

}

// spamControl contains functions for managing
// spam control filters of Google Analytics accounts.
type spamControl struct {
	analyticsAPI *api.API
}

// Update the spam control filters.
func (cli *spamControl) Update() error {
	return fmt.Errorf("No implemented")
}

// Remove all spam control filters.
func (cli *spamControl) Remove() error {
	return fmt.Errorf("No implemented")
}

// Status displays the spam control status.
func (cli *spamControl) Status() error {

	// get all available accounts
	accounts, accountsError := cli.analyticsAPI.GetAccounts()
	if accountsError != nil {
		return accountsError
	}

	for _, account := range accounts {

		// get all filters for account
		filters, filtersError := cli.analyticsAPI.GetFilters(account.ID)
		if filtersError != nil {
			return filtersError
		}

		for _, filter := range filters {
			log.Printf("%#v\n", filter)
		}

	}

	return nil

}
