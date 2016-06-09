package spamcontrol

import (
	"sort"

	"github.com/andreaskoch/ga-spam-control/api"
)

// NewSpamDomainProviderFactory creates a new spam-domain provider factory instance.
func NewSpamDomainProviderFactory(analyticsAPI api.AnalyticsAPI, spamDetector SpamDetector) *SpamDomainProviderFactory {
	spamAnalysis := &dynamicSpamAnalysis{
		analyticsDataProvider: &remoteAnalyticsDataProvider{
			analyticsAPI: analyticsAPI,
		},
		spamDetector: spamDetector,
	}

	singleAccountAnalysis := &SingleAccountAnalysis{spamAnalysis}

	return &SpamDomainProviderFactory{analyticsAPI, singleAccountAnalysis}
}

// SpamDomainProviderFactory allows to get a list of multiple
// spamDomainProvider instances based on the accounts that are accessible by
// the current user via the given analytics API instance.
type SpamDomainProviderFactory struct {
	analyticsAPI          api.AnalyticsAPI
	singleAccountAnalysis *SingleAccountAnalysis
}

// GetSpamDomainProviders returns a slice of spamDomainProvider instances.
func (factory *SpamDomainProviderFactory) GetSpamDomainProviders() ([]SpamDomainProvider, error) {
	accounts, err := factory.analyticsAPI.GetAccounts()
	if err != nil {
		return nil, err
	}

	var accountIDs []string
	for _, account := range accounts {
		accountIDs = append(accountIDs, account.ID)
	}

	ddoforgGhostSpamListProvider := &staticSpamDomains{"https://raw.githubusercontent.com/ddofborg/analytics-ghost-spam-list/master/adwordsrobot.com-spam-list.txt"}
	stevieRayApacheNginxSpamListProvider := &staticSpamDomains{"https://raw.githubusercontent.com/Stevie-Ray/apache-nginx-referral-spam-blacklist/master/generator/domains.txt"}
	piwikSpamListProvider := &staticSpamDomains{"https://raw.githubusercontent.com/piwik/referrer-spam-blacklist/master/spammers.txt"}

	multiAccountAnalysis := &MultiAccountAnalysis{accountIDs, factory.singleAccountAnalysis}

	return []SpamDomainProvider{
		ddoforgGhostSpamListProvider,
		stevieRayApacheNginxSpamListProvider,
		piwikSpamListProvider,
		multiAccountAnalysis,
	}, nil
}

// MultiAccountAnalysis is a spamDomainProvider that aggregates the results
// of multiple analytics accounts.
type MultiAccountAnalysis struct {
	accountIDs            []string
	singleAccountAnalysis *SingleAccountAnalysis
}

// GetSpamDomains returns a list of referrer spam domain names that is extracted
// from the latest analytics data and rated by a spam-detector for multiple
// analytics accounts.
func (provider *MultiAccountAnalysis) GetSpamDomains() ([]string, error) {

	var spamDomains []string
	for _, accountID := range provider.accountIDs {
		domains, err := provider.singleAccountAnalysis.GetSpamDomains(accountID)
		if err != nil {
			return nil, err
		}

		spamDomains = append(spamDomains, domains...)
	}

	spamDomains = removeDuplicatesFromList(spamDomains)
	sort.Strings(spamDomains)

	return spamDomains, nil
}

// The SingleAccountAnalysis struct provides an
// up-to-date list of spam domain for a given analytics account.
type SingleAccountAnalysis struct {
	spamAnalysis spamAnalysis
}

// GetSpamDomains returns a list of referrer spam domain names that is extracted
// from the latest analytics data and rated by a spam-detector.
func (provider *SingleAccountAnalysis) GetSpamDomains(accountID string) ([]string, error) {

	analysis, err := provider.spamAnalysis.GetSpamAnalysis(accountID, 30, 0.75)
	if err != nil {
		return nil, err
	}

	var domainNames []string
	for _, spamDomain := range analysis.SpamDomains {
		domainNames = append(domainNames, spamDomain.DomainName)
	}

	sort.Strings(domainNames)

	return domainNames, nil
}
