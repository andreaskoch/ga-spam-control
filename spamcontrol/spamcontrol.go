// Package spamcontrol contains all business logic for spam-filter manipulations.
package spamcontrol

import (
	"fmt"
	"regexp"

	"github.com/andreaskoch/ga-spam-control/api"
)

// The SpamController interface provides functions for displaying, updating and removing referrer spam controls.
type SpamController interface {

	// Remove the referrer spam controls from the account with the given accountID.
	// Returns an error if the removal failed.
	Remove(accountID string) error

	// Analyze checks the given account for referrer spam and returns the result
	// of the analysis as a view model. Returns an error if the analysis failed.
	DetectSpam(accountID string, numberOfDaysToLookBack int) (AnalysisResult, error)

	// UpdateSpamDomains updates the referrer spam domain list.
	UpdateSpamDomains() (UpdateResult, error)

	// ListSpamDomains returns a list of all known spam domains
	ListSpamDomains() ([]string, error)

	// Status collects the current spam-control status of all accessible
	// analytics accounts. It returns the a StateOverview model with the Status
	// of all accounts and an overall Status. If the status cannot be determined
	// an error will be returned.
	GlobalStatus() (StateOverview, error)

	// AccountStatus returns the current spam-control status of the account
	// with the given account ID. Returns an error if the status cannot be
	// determined.
	AccountStatus(accountID string) (InstallationStatus, error)

	// Update the referrer spam controls for the account with the given accountID.
	// Returns an error if the update failed.
	Update(accountID string) error
}

// New creates a new spam control instance.
func New(analyticsAPI api.AnalyticsAPI, spamDomainProvider SpamDomainProvider, communitySpamRepository *CommunitySpamDomainRepository, privateSpamRepository *PrivateSpamDomainRepository) *SpamControl {

	accountProvider := remoteAccountProvider{analyticsAPI}

	spamDetector := &interactiveSpamDetector{
		analyticsDataProvider: &remoteAnalyticsDataProvider{
			analyticsAPI: analyticsAPI,
		},
		spamDomainProvider: spamDomainProvider,
	}

	filterNameProvider := &spamFilterNameProvider{"Referrer Spam Block"}

	filterFactory := &googleAnalyticsFilterFactory{
		filterNameProvider:   filterNameProvider,
		filterValueMaxLength: 255,
	}

	filterProvider := &remoteFilterProvider{
		analyticsAPI: analyticsAPI,

		spamDomainProvider: spamDomainProvider,

		filterNameProvider: filterNameProvider,
		filterFactory:      filterFactory,
	}

	return &SpamControl{
		accountProvider: accountProvider,
		filterProvider:  filterProvider,
		spamDetector:    spamDetector,

		spamDomainProvider: spamDomainProvider,

		communitySpamRepository: communitySpamRepository,
		privateSpamRepository:   privateSpamRepository,
	}
}

// The SpamControl type provides functions for
// managing Google Analtics spam filters.
type SpamControl struct {
	accountProvider accountProvider
	filterProvider  filterProvider
	spamDetector    spamDetector

	spamDomainProvider SpamDomainProvider

	// repositories
	communitySpamRepository *CommunitySpamDomainRepository
	privateSpamRepository   *PrivateSpamDomainRepository
}

// Remove the referrer spam controls from the account with the given accountID.
// Returns an error if the removal failed.
func (spamControl *SpamControl) Remove(accountID string) error {

	// get the requested account
	account, accountError := spamControl.accountProvider.GetAccount(accountID)
	if accountError != nil {
		return accountError
	}

	// get all filters for account
	filters, filtersError := spamControl.filterProvider.GetExistingFilters(account.ID)
	if filtersError != nil {
		return filtersError
	}

	for _, filter := range filters {
		if err := spamControl.filterProvider.RemoveFilter(account.ID, filter.ID); err != nil {
			return err
		}
	}

	return nil
}

// ListSpamDomains returns a list of all known spam domains
func (spamControl *SpamControl) ListSpamDomains() ([]string, error) {
	return spamControl.spamDomainProvider.GetSpamDomains()
}

// UpdateSpamDomains updates the referrer spam domain list.
func (spamControl *SpamControl) UpdateSpamDomains() (UpdateResult, error) {
	unchanged, added, removed, err := spamControl.communitySpamRepository.UpdateSpamDomains()
	if err != nil {
		return UpdateResult{}, err
	}

	return createUpdateResultModel(unchanged, added, removed)
}

// DetectSpam checks the given account for referrer spam.
// Returns an error if the analysis failed.
func (spamControl *SpamControl) DetectSpam(accountID string, numberOfDaysToLookBack int) (AnalysisResult, error) {

	// get the requested account
	account, accountError := spamControl.accountProvider.GetAccount(accountID)
	if accountError != nil {
		return AnalysisResult{}, accountError
	}

	if numberOfDaysToLookBack < 1 {
		return AnalysisResult{}, fmt.Errorf("The specified number of days to look back cannot be below 1")
	}

	// find new referrer spam domains
	newSpamDomainNames, err := spamControl.spamDetector.DetectSpam(account.ID, numberOfDaysToLookBack)
	if err != nil {
		return AnalysisResult{}, err
	}

	// assemble the viewmodel from the result
	result := AnalysisResult{
		AccountID: accountID,
	}

	for _, domainName := range newSpamDomainNames {

		result.Domains = append(result.Domains, Domain{
			DomainName: domainName,
			IsSpam:     true,
		})

	}

	// add the new domain names to the repository
	addError := spamControl.privateSpamRepository.AddDomains(newSpamDomainNames)
	if addError != nil {
		return AnalysisResult{}, err
	}

	return result, nil
}

// GlobalStatus collects the current spam-control status of all accessible
// analytics accounts. It returns the a StateOverview model with the Status
// of all accounts and an overall Status. If the status cannot be determined
// an error will be returned.
func (spamControl *SpamControl) GlobalStatus() (StateOverview, error) {
	// get all available accounts
	accounts, accountsError := spamControl.accountProvider.GetAccounts()
	if accountsError != nil {
		return StateOverview{}, accountsError
	}

	knownSpamDomains, spamDomainsError := spamControl.spamDomainProvider.GetSpamDomains()
	if spamDomainsError != nil {
		return StateOverview{}, spamDomainsError
	}

	overviewModel := &StateOverview{
		Accounts:         make([]AccountStatus, 0),
		KnownSpamDomains: len(knownSpamDomains),
	}

	// get the status for each account
	for _, account := range accounts {

		accountStatus, accountStatusError := spamControl.AccountStatus(account.ID)
		if accountStatusError != nil {
			return StateOverview{}, accountStatusError
		}

		accountStatusModel := AccountStatus{
			AccountID:   account.ID,
			AccountName: account.Name,
			Status:      accountStatus,
		}

		overviewModel.Accounts = append(overviewModel.Accounts, accountStatusModel)
	}

	return *overviewModel, nil
}

// AccountStatus returns the current spam-control status of the account
// with the given account ID. Returns an error if the status cannot be
// determined.
func (spamControl *SpamControl) AccountStatus(accountID string) (InstallationStatus, error) {

	// get the requested account
	account, accountError := spamControl.accountProvider.GetAccount(accountID)
	if accountError != nil {
		return InstallationStatus{}, accountError
	}

	// get the existing filters
	existingFilters, existingFilterError := spamControl.filterProvider.GetExistingFilters(account.ID)
	if existingFilterError != nil {
		return InstallationStatus{}, existingFilterError
	}

	// get the latest referrer spam domain names
	domainNames, spamDomainProviderError := spamControl.spamDomainProvider.GetSpamDomains()
	if spamDomainProviderError != nil {
		return InstallationStatus{}, spamDomainProviderError
	}

	// test the filters
	domainsCovered := make(map[string]int)
	for _, filter := range existingFilters {
		expressionRegex := regexp.MustCompile(filter.ExcludeDetails.ExpressionValue)

		for _, domain := range domainNames {
			if isMatch := expressionRegex.MatchString(domain); !isMatch {
				continue
			}

			domainsCovered[domain]++
		}

	}

	return InstallationStatus{
		TotalDomains:   len(domainNames),
		DomainsCovered: len(domainsCovered),
	}, nil
}

// Update the referrer spam controls for the account with the given accountID.
// Returns an error if the update failed.
func (spamControl *SpamControl) Update(accountID string) error {

	// get the requested account
	account, accountError := spamControl.accountProvider.GetAccount(accountID)
	if accountError != nil {
		return accountError
	}

	return spamControl.filterProvider.UpdateFilters(account.ID)

}

// createUpdateResultModel creates an UpdateResult model from the given lists
// of unchaged, added and removed spam domain names.
func createUpdateResultModel(unchanged, added, removed []string) (UpdateResult, error) {
	var domainUpdates []DomainUpdate

	// unchanged
	for _, domain := range unchanged {
		domainUpdates = append(domainUpdates, DomainUpdate{
			UpdateType: Unchanged,
			Domainname: domain,
		})
	}

	// added
	for _, domain := range added {
		domainUpdates = append(domainUpdates, DomainUpdate{
			UpdateType: Added,
			Domainname: domain,
		})
	}

	// removed
	for _, domain := range removed {
		domainUpdates = append(domainUpdates, DomainUpdate{
			UpdateType: Removed,
			Domainname: domain,
		})
	}

	// return an error if the list if empty
	if len(domainUpdates) == 0 {
		return UpdateResult{}, fmt.Errorf("No domains received.")
	}

	// sort the list
	SortDomainUpdateBy(domainUpdatesByName).Sort(domainUpdates)

	return UpdateResult{
		Statistics: DomainUpdateStatistics{
			Unchanged: len(unchanged),
			Added:     len(added),
			Removed:   len(removed),
			Total:     len(domainUpdates),
		},
		Domains: domainUpdates,
	}, nil
}
