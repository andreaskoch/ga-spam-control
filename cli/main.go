// Package cli provides the commandline interface for the spamcontrol package.
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

// configFolderName contains the name of the folder that
// is used for storing credentials and other configuration files.
const configFolderName = ".ga-spam-control"

// tokenStoreFilePath contains the file path to the file which holds the oAuth
// credentials for the Google Analytics API.
var tokenStoreFilePath string

// communitySpamListFilePath contains the file path to the file which stores
// a local copy of the aggregated referrer spam lists maintained by the community.
var communitySpamListFilePath string

// personalSpamListFilePath contains the file path to the file which stores
// your personal referrer spam list.
var personalSpamListFilePath string

func init() {
	homeDirPath, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot determine the current users home direcotry. Error: %s", err)
		os.Exit(1)
	}

	tokenStoreFilePath = filepath.Join(homeDirPath, configFolderName, "credentials.json")
	communitySpamListFilePath = filepath.Join(homeDirPath, configFolderName, "spam-domains", "community.txt")
	personalSpamListFilePath = filepath.Join(homeDirPath, configFolderName, "spam-domains", "personal.txt")
}

func main() {
	handleCommandlineArguments(os.Args[1:])
}

// handleCommandlineArguments parses the given arguments
// and performs the selected action.
func handleCommandlineArguments(args []string) {
	app := kingpin.New("ga-spam-control", "Command-line utility for keeping your Google Analytics referrer spam filters up-to-date")
	app.Version("0.6.0")

	// filters
	filtersCommand := app.Command("filters", "Manage spam filters")

	statusCommand := filtersCommand.Command("status", "Show the spam-control status of your accounts")
	statusCommandAccountID := statusCommand.Arg("accountID", "Google Analytics account ID").String()
	statusCommandQuiet := statusCommand.Flag("quiet", "Display spam-protection status in a parsable format").Short('q').Bool()

	updateFiltersCommand := filtersCommand.Command("update", "Update the spam filters for the given account")
	updateFiltersAccountID := updateFiltersCommand.Arg("accountID", "Google Analytics account ID").Required().String()

	removeFiltersCommand := filtersCommand.Command("remove", "Remove all spam filters from an account")
	removeFiltersAccountID := removeFiltersCommand.Arg("accountID", "Google Analytics account ID").Required().String()

	// domains
	domainsCommand := app.Command("domains", "Manage spam domains")
	listDomainsCommand := domainsCommand.Command("list", "List all currently known spam domains")

	updateSpamDomainsCommand := domainsCommand.Command("update", fmt.Sprintf("Update your list of known referrer spam domains (%q)", communitySpamListFilePath))
	updateSpamDomainsQuiet := updateSpamDomainsCommand.Flag("quiet", "Display results in a parsable format").Short('q').Bool()

	findSpamDomains := domainsCommand.Command("find", fmt.Sprintf("Find new referrer spam domains in your analytics data and write them to your private referrer spam list (%q)", personalSpamListFilePath))
	findSpamDomainsAccountID := findSpamDomains.Arg("accountID", "Google Analytics account ID").Required().String()
	findSpamDomainsNumberOfDays := findSpamDomains.Arg("days", "The number of days to look back").Default("90").Int()
	findSpamDomainsQuiet := findSpamDomains.Flag("quiet", "Display results in a parsable format").Short('q').Bool()

	command, err := app.Parse(args)
	if err != nil {
		app.Fatalf("%s", err.Error())
	}

	switch command {

	case statusCommand.FullCommand():
		// Display status
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		statusError := cli.Status(*statusCommandAccountID, *statusCommandQuiet)
		if statusError != nil {
			app.Fatalf("%s", statusError.Error())
		}

	case updateFiltersCommand.FullCommand():
		// Update filters
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		updateError := cli.Update(*updateFiltersAccountID)
		if updateError != nil {
			app.Fatalf("%s", updateError.Error())
		}

	case removeFiltersCommand.FullCommand():
		// Remove spam control
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		removeError := cli.Remove(*removeFiltersAccountID)
		if removeError != nil {
			app.Fatalf("%s", removeError.Error())
		}

	case updateSpamDomainsCommand.FullCommand():
		// update spam domains
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		updateSpamDomainsError := cli.UpdateSpamDomains(*updateSpamDomainsQuiet)
		if updateSpamDomainsError != nil {
			app.Fatalf("%s", updateSpamDomainsError.Error())
		}

	case listDomainsCommand.FullCommand():
		// list spam domains
		cli, err := newCLI()
		if err != nil {
			app.Fatalf("%s", err.Error())
		}

		findSpamDomainsError := cli.ListSpamDomains()
		if findSpamDomainsError != nil {
			app.Fatalf("%s", findSpamDomainsError.Error())
		}

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

	}
}

// newCLI creates a new spam control instance.
func newCLI() (*cli, error) {

	// create a token store
	tokenStore := token.NewTokenStore(tokenStoreFilePath)

	// create a new analytis API instance
	googleAnalyticsClientID := "821429244906-8aki1tiaov6g2o7lr7elp41435adk9ge.apps.googleusercontent.com"
	googleAnalyticsClientSecret := "_WxLj0SpQ8HxqmOEyYDUTFzW"
	analyticsAPI, apiError := api.New(tokenStore, googleAnalyticsClientID, googleAnalyticsClientSecret)
	if apiError != nil {
		return nil, apiError
	}

	// community spam domain repository
	remoteSpamDomainURLs := []string{
		"https://raw.githubusercontent.com/ddofborg/analytics-ghost-spam-list/master/adwordsrobot.com-spam-list.txt",
		"https://raw.githubusercontent.com/Stevie-Ray/apache-nginx-referral-spam-blacklist/master/generator/domains.txt",
		"https://raw.githubusercontent.com/piwik/referrer-spam-blacklist/master/spammers.txt",
	}
	remoteProviders := getRemoteSpamDomainProviders(remoteSpamDomainURLs)

	communitySpamListRepository := spamcontrol.NewCommunitySpamDomainRepository(communitySpamListFilePath, remoteProviders)

	// personal spam domain repository
	localProvider := spamcontrol.NewLocalSpamDomainProvider(personalSpamListFilePath)
	privateSpamListRepository := spamcontrol.NewPrivateSpamDomainRepository(personalSpamListFilePath, localProvider)

	// combined private and community spam domain repository
	combinedSpamDomainProvider := spamcontrol.NewAggregateProvider([]spamcontrol.SpamDomainProvider{
		communitySpamListRepository,
		privateSpamListRepository,
	})

	// create a spam control instance
	spamControl := spamcontrol.New(
		analyticsAPI,
		combinedSpamDomainProvider,
		communitySpamListRepository,
		privateSpamListRepository)

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

func getRemoteSpamDomainProviders(urls []string) []spamcontrol.SpamDomainProvider {
	var providers []spamcontrol.SpamDomainProvider
	for _, url := range urls {
		providers = append(providers, spamcontrol.NewRemoteSpamDomainProvider(url))
	}
	return providers
}
