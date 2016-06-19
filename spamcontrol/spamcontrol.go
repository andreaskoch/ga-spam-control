// Package spamcontrol contains all business logic for spam-filter manipulations.
package spamcontrol

import (
	"fmt"

	"github.com/andreaskoch/ga-spam-control/api"
	"github.com/andreaskoch/ga-spam-control/spamcontrol/status"
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
	if numberOfDaysToLookBack < 1 {
		return AnalysisResult{}, fmt.Errorf("The specified number of days to look back cannot be below 1")
	}

	// find new referrer spam domains
	newSpamDomainNames, err := spamControl.spamDetector.DetectSpam(accountID, numberOfDaysToLookBack)
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

		accountStatus, accountStatusError := spamControl.filterProvider.GetAccountStatus(account.ID)
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

	// get the accounts' status
	installationStatus, accountStatusError := spamControl.filterProvider.GetAccountStatus(account.ID)
	if accountStatusError != nil {
		return InstallationStatus{}, accountStatusError
	}

	return installationStatus, nil
}

// Update the referrer spam controls for the account with the given accountID.
// Returns an error if the update failed.
func (spamControl *SpamControl) Update(accountID string) error {

	// get the requested account
	account, accountError := spamControl.accountProvider.GetAccount(accountID)
	if accountError != nil {
		return accountError
	}

	filterStatuses, filterStatusError := spamControl.filterProvider.GetFilterStatuses(account.ID)
	if filterStatusError != nil {
		return filterStatusError
	}

	for _, filterStatus := range filterStatuses {

		// ignore up-to-date filters
		if filterStatus.Status() == status.UpToDate {
			continue
		}

		// remove obsolete filters
		if filterStatus.Status() == status.Obsolete {
			fmt.Printf("Removing filter %q"+NewLineSequence, filterStatus.Filter().Name)
			removeError := spamControl.filterProvider.RemoveFilter(account.ID, filterStatus.Filter().ID)
			if removeError != nil {
				return removeError
			}

			continue
		}

		// update outdated filters
		if filterStatus.Status() == status.Outdated {
			fmt.Printf("Updating filter %q"+NewLineSequence, filterStatus.Filter().Name)
			_, updateError := spamControl.filterProvider.UpdateFilter(account.ID, filterStatus.Filter().ID, filterStatus.Filter())
			if updateError != nil {
				return updateError
			}

			continue
		}

		// create new filters
		if filterStatus.Status() == status.NotInstalled {
			fmt.Printf("Creating filter %q"+NewLineSequence, filterStatus.Filter().Name)
			_, createError := spamControl.filterProvider.CreateFilter(account.ID, filterStatus.Filter())
			if createError != nil {
				return createError
			}

			continue
		}

		return fmt.Errorf("Cannot update filter %q. Status %q cannot be handled.", filterStatus.Filter().Name, filterStatus.Status())
	}

	return nil
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
		return UpdateResult{}, fmt.Errorf("Something is wrong. No domains received.")
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
