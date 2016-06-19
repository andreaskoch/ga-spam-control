package spamcontrol

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/andreaskoch/ga-spam-control/api"
)

type spamDetector interface {
	// DetectSpam returns a list of spam domain names for a given account ID.
	DetectSpam(accountID string, numberOfDays int) ([]string, error)
}

// interactiveSpamDetector performs an interactive referrer spam analysis on the analytics data for the given acount.
type interactiveSpamDetector struct {
	analyticsDataProvider analyticsDataProvider
	spamDomainProvider    SpamDomainProvider
	knownSpamDomains      map[string]string
}

// DetectSpam returns a list of spam domain names for a given account ID.
func (spamDetector *interactiveSpamDetector) DetectSpam(accountID string, numberOfDays int) ([]string, error) {

	// get the analytics data for the given account
	analyticsData, analyticsDataError := spamDetector.analyticsDataProvider.GetAnalyticsData(accountID, numberOfDays)
	if analyticsDataError != nil {
		return nil, analyticsDataError
	}

	initialReferralData := spamDetector.getReferralData(analyticsData)
	newSpamDomainNames, reviewError := spamDetector.reviewReferralData(initialReferralData)
	if reviewError != nil {
		return nil, reviewError
	}

	return newSpamDomainNames, nil
}

func (spamDetector *interactiveSpamDetector) reviewReferralData(reviewData ReferralData) ([]string, error) {

	// write the domain names to a temp file
	tmpDir := os.TempDir()
	reviewFile, tmpFileErr := ioutil.TempFile(tmpDir, "referrer-spam-domain.review")
	if tmpFileErr != nil {
		return nil, fmt.Errorf("Error %s while creating tempFile", tmpFileErr)
	}

	fmt.Fprintf(reviewFile, "# %s"+NewLineSequence, "Review referrer domain names.")
	fmt.Fprintf(reviewFile, "# %s"+NewLineSequence, "Uncomment all lines which you consider to be referrer spam.")
	fmt.Fprintf(reviewFile, "%s"+NewLineSequence, "")

	reviewData.ForEach(func(index int, domain ReferralDomain, analyticsData api.AnalyticsData) {

		if !domain.IsSpam() {
			fmt.Fprintf(reviewFile, "#%s"+NewLineSequence, domain.Name())
		}

	})

	filePath := reviewFile.Name()
	reviewFile.Close()

	// open the temp file in a text editor
	editorPath, editorPathLookupError := getTextEditorPath()
	if editorPathLookupError != nil {
		return nil, editorPathLookupError
	}

	textEditorCommand := exec.Command(editorPath, filePath)
	textEditorCommand.Stdin = os.Stdin
	textEditorCommand.Stdout = os.Stdout

	if startError := textEditorCommand.Start(); startError != nil {
		return nil, fmt.Errorf("Start failed: %s", startError)
	}

	if editorCommandError := textEditorCommand.Wait(); editorCommandError != nil {
		return nil, fmt.Errorf("Command finished with error: %v", editorCommandError)
	}

	// read the results
	reviewFile, readError := os.Open(filePath)
	if readError != nil {
		return nil, fmt.Errorf("Unable to open %q. Error: %s", filePath, readError)
	}

	defer reviewFile.Close()

	var spamDomainNames []string
	scanner := bufio.NewScanner(reviewFile)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())

		// ignore empty lines
		if domain == "" {
			continue
		}

		// ignore commented lines
		if strings.HasPrefix(domain, "#") {
			continue
		}

		spamDomainNames = append(spamDomainNames, domain)
	}

	unique(spamDomainNames)
	sort.Strings(spamDomainNames)

	return spamDomainNames, nil
}

// getReferralData extracts ReferralData from the given analytics data.
func (spamDetector *interactiveSpamDetector) getReferralData(analyticsData api.AnalyticsData) ReferralData {

	var domainNames []string
	data := make(map[string]api.AnalyticsData)

	for _, row := range analyticsData {

		// ignore non-referrer entries
		if row.Medium != "referral" {
			continue
		}

		// ignore empty values
		if row.Source == "(direct)" || row.Source == "" {
			continue
		}

		// the source contains the domain name
		domainName := row.Source

		domainNames = append(domainNames, domainName)
		data[domainName] = append(data[domainName], row)
	}

	// removed duplicates
	domainNames = unique(domainNames)

	// sort the list alphabetically
	sort.Strings(domainNames)

	// pre-rate the domain names
	var domains []ReferralDomain
	for _, domainName := range domainNames {
		domains = append(domains, ReferralDomain{
			domainName,
			spamDetector.isKnownSpamDomain(domainName),
		})
	}

	return ReferralData{domains, data}
}

// isKnownSpamDomain checks if the given domain name is a known referrer spam domain.
func (spamDetector *interactiveSpamDetector) isKnownSpamDomain(domainName string) bool {

	// initialize the map (once!)
	if spamDetector.knownSpamDomains == nil {
		spamDetector.knownSpamDomains = make(map[string]string)

		knownDomains, err := spamDetector.spamDomainProvider.GetSpamDomains()
		if err != nil {
			panic(err)
		}

		for _, knownSpamDomain := range knownDomains {
			spamDetector.knownSpamDomains[knownSpamDomain] = knownSpamDomain
		}
	}

	if _, exists := spamDetector.knownSpamDomains[domainName]; exists {
		return true
	}

	return false
}

// A ReferralDomain represents a domain name and a spam state.
type ReferralDomain struct {
	domainName string
	isSpam     bool
}

// Name returns the domain name of the referral domain.
func (referralDomain ReferralDomain) Name() string {
	return referralDomain.domainName
}

// IsSpam returns a flag indicating whether the current referral domain is spam or not.
func (referralDomain ReferralDomain) IsSpam() bool {
	return referralDomain.isSpam
}

// ReferralData contains the analytics data for a single referrer domain.
type ReferralData struct {
	domains []ReferralDomain
	data    map[string]api.AnalyticsData
}

// Data returns the list of referral domain names.
func (referralData ReferralData) Data(domainName string) (api.AnalyticsData, bool) {
	if data, exists := referralData.data[domainName]; exists {
		return data, true
	}

	return api.AnalyticsData{}, false
}

// Domains returns the list of referral domain names.
func (referralData ReferralData) Domains() []ReferralDomain {
	return referralData.domains
}

// ForEach iterates over each domain name in this ReferralData object and passes
// the domain name and the analytics data for the domain name to the given expression.
func (referralData ReferralData) ForEach(expression func(index int, domain ReferralDomain, analyticsData api.AnalyticsData)) {
	for index, domain := range referralData.domains {
		expression(index, domain, referralData.data[domain.Name()])
	}
}

// getTextEditorPath returns the path of an text editor or
// returns an error if no text editor was found.
func getTextEditorPath() (string, error) {
	textEditor := DefaultEditor
	editorPath, editorPathLookupError := exec.LookPath(textEditor)
	if editorPathLookupError != nil {
		return "", fmt.Errorf("Error %s while looking up for %s. Error: %s", editorPath, textEditor, editorPathLookupError)
	}

	return editorPath, nil
}
