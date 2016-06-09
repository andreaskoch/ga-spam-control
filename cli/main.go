// Package cli provides the commandline interface for the spamcontrol package.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	app := kingpin.New("ga-spam-control", "Command-line utility for keeping your Google Analytics referrer spam filters up-to-date")
	app.Version("0.4.0")

	status := app.Command("show-status", "Display the spam-control status of your accounts")
	statusAccountID := status.Arg("accountID", "Google Analytics account ID").String()
	statusQuiet := status.Flag("quiet", "Display spam-protection status in a parsable format").Short('q').Bool()

	update := app.Command("update-filters", "Update the spam filters for the given account")
	updateAccountID := update.Arg("accountID", "Google Analytics account ID").Required().String()

	remove := app.Command("remove-filters", "Remove all spam filters from an account")
	removeAccountID := remove.Arg("accountID", "Google Analytics account ID").Required().String()

	listSpamDomains := app.Command("list-spam-domains", "List all currently known spam domains")

	updateSpamDomains := app.Command("update-spam-domains", "Update the spam domain list")
	updateSpamDomainsQuiet := updateSpamDomains.Flag("quiet", "Display the analyis results in a parsable format").Short('q').Bool()

	findSpamDomains := app.Command("find-spam-domains", "Find new referrer spam domains in your analytics data")
	findSpamDomainsAccountID := findSpamDomains.Arg("accountID", "Google Analytics account ID").Required().String()
	findSpamDomainsNumberOfDays := findSpamDomains.Arg("days", "The number of days to look back").Default("3").Int()
	findSpamDomainsQuiet := findSpamDomains.Flag("quiet", "Display the analyis results in a parsable format").Short('q').Bool()

	getTrainingData := app.Command("get-training-data", "Get training data for the given account")
	getTrainingDataAccountID := getTrainingData.Arg("accountID", "Google Analytics account ID").Required().String()
	getTrainingDataNumberOfDays := getTrainingData.Arg("days", "The number of days to look back").Default("3").Int()

	switch kingpin.MustParse(app.Parse(args)) {

	case status.FullCommand():
		// Display status
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		statusError := cli.Status(*statusAccountID, *statusQuiet)
		if statusError != nil {
			app.Fatalf("%s", statusError.Error())
		}

		os.Exit(0)

	case update.FullCommand():
		// Update filters
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		updateError := cli.Update(*updateAccountID)
		if updateError != nil {
			app.Fatalf("%s", updateError.Error())
		}

		os.Exit(0)

	case remove.FullCommand():
		// Remove spam control
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		removeError := cli.Remove(*removeAccountID)
		if removeError != nil {
			app.Fatalf("%s", removeError.Error())
		}

		os.Exit(0)

	case updateSpamDomains.FullCommand():
		// update spam domains
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		updateSpamDomainsError := cli.UpdateSpamDomains(*updateSpamDomainsQuiet)
		if updateSpamDomainsError != nil {
			app.Fatalf("%s", updateSpamDomainsError.Error())
		}

		os.Exit(0)

	case listSpamDomains.FullCommand():
		// list spam domains
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		findSpamDomainsError := cli.ListSpamDomains()
		if findSpamDomainsError != nil {
			app.Fatalf("%s", findSpamDomainsError.Error())
		}

		os.Exit(0)

	case findSpamDomains.FullCommand():
		// find spam
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		findSpamDomainsError := cli.DetectSpam(*findSpamDomainsAccountID, *findSpamDomainsNumberOfDays, *findSpamDomainsQuiet)
		if findSpamDomainsError != nil {
			app.Fatalf("%s", findSpamDomainsError.Error())
		}

		os.Exit(0)

	case getTrainingData.FullCommand():
		// get training data
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		getTrainingDataError := cli.GetTrainingData(*getTrainingDataAccountID, *getTrainingDataNumberOfDays)
		if getTrainingDataError != nil {
			app.Fatalf("%s", getTrainingDataError.Error())
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

	tokenStoreFilePath := filepath.Join(homeDirPath, ".ga-spam-control", "credentials")
	tokenStore := token.NewTokenStore(tokenStoreFilePath)

	// create a new analytis API instance
	googleAnalyticsClientID := "821429244906-8aki1tiaov6g2o7lr7elp41435adk9ge.apps.googleusercontent.com"
	googleAnalyticsClientSecret := "_WxLj0SpQ8HxqmOEyYDUTFzW"
	analyticsAPI, apiError := api.New(tokenStore, googleAnalyticsClientID, googleAnalyticsClientSecret)
	if apiError != nil {
		return nil, apiError
	}

	spamDetector := spamcontrol.NewDetector()
	domainProviderFactory := spamcontrol.NewSpamDomainProviderFactory(analyticsAPI, spamDetector)

	spamDomainRepositoryFilePath := filepath.Join(homeDirPath, ".ga-spam-control", "domains")
	domainRepository := spamcontrol.NewSpamDomainRepository(spamDomainRepositoryFilePath, domainProviderFactory)

	// create a spam control instance
	spamControl := spamcontrol.New(analyticsAPI, spamDetector, domainRepository)

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

// UpdateSpamDomains updates the referrer spam domain lists.
func (cli *cli) UpdateSpamDomains(quiet bool) error {
	updateResult, err := cli.spamControl.UpdateSpamDomains()
	if err != nil {
		return err
	}

	// select the display template
	templateText := templates.PrettyUpdate
	if quiet {
		// the "quiet" template contains no surplus texts
		// and should be easier to parse by tools like awk
		templateText = templates.QuietUpdate
	}

	updateResultTemplate, parseError := template.New("UpdateResult").Parse(templateText)
	if parseError != nil {
		return parseError
	}

	renderError := updateResultTemplate.Execute(os.Stdout, updateResult)
	if renderError != nil {
		return renderError
	}

	return nil
}

// ListSpamDomains prints a list of all known referrer spam domains.
func (cli *cli) ListSpamDomains() error {
	domains, err := cli.spamControl.ListSpamDomains()
	if err != nil {
		return err
	}

	for _, domain := range domains {
		fmt.Println(domain)
	}

	return nil
}

// DetectSpam checks the given account for referrer-spam.
func (cli *cli) DetectSpam(accountID string, numberOfDaysToLookBack int, quiet bool) error {
	analysisResultViewModel, err := cli.spamControl.DetectSpam(accountID, numberOfDaysToLookBack)
	if err != nil {
		return err
	}

	// select the display template
	templateText := templates.PrettyAnalysis
	if quiet {
		// the "quiet" template contains no surplus texts
		// and should be easier to parse by tools like awk
		templateText = templates.QuietAnalysis
	}

	analysisTemplate, parseError := template.New("Analysis").Parse(templateText)
	if parseError != nil {
		return parseError
	}

	renderError := analysisTemplate.Execute(os.Stdout, analysisResultViewModel)
	if renderError != nil {
		return renderError
	}

	return nil
}

func (cli *cli) GetTrainingData(accountID string, numberOfDaysToLookBack int) error {
	trainingData, err := cli.spamControl.GetTrainingData(accountID, numberOfDaysToLookBack)
	if err != nil {
		return err
	}

	separator := ","
	fmt.Fprintf(os.Stdout, "%s\n", strings.Join(trainingData.ColumnNames, separator))
	for _, row := range trainingData.Rows {
		fmt.Fprintf(os.Stdout, "%s\n", strings.Join(row, separator))
	}

	return nil
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
