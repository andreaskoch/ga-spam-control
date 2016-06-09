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

	// GetTrainingData returns a set of training data for the given accountID.
	// Returns an error if the training data could not be fetched.
	GetTrainingData(accountID string, numberOfDaysToLookBack int) (MachineLearningModel, error)

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
func New(analyticsAPI api.AnalyticsAPI, spamDetector SpamDetector, spamRepository SpamDomainRepository) *SpamControl {

	accountProvider := remoteAccountProvider{analyticsAPI}

	spamAnalysis := &dynamicSpamAnalysis{
		analyticsDataProvider: &remoteAnalyticsDataProvider{
			analyticsAPI: analyticsAPI,
		},
		spamDetector: spamDetector,
	}

	filterNameProvider := &spamFilterNameProvider{"Referrer Spam Block"}

	filterFactory := &googleAnalyticsFilterFactory{
		filterNameProvider:   filterNameProvider,
		filterValueMaxLength: 255,
	}

	filterProvider := &remoteFilterProvider{
		analyticsAPI: analyticsAPI,

		spamRepository: spamRepository,

		filterNameProvider: filterNameProvider,
		filterFactory:      filterFactory,
	}

	return &SpamControl{
		accountProvider: accountProvider,
		filterProvider:  filterProvider,
		spamAnalysis:    spamAnalysis,
		spamRepository:  spamRepository,
	}
}

// The SpamControl type provides functions for
// managing Google Analtics spam filters.
type SpamControl struct {
	accountProvider accountProvider
	filterProvider  filterProvider
	spamAnalysis    spamAnalysis
	spamRepository  SpamDomainRepository
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

// UpdateSpamDomains updates the referrer spam domain list.
func (spamControl *SpamControl) UpdateSpamDomains() (UpdateResult, error) {
	unchanged, added, removed, err := spamControl.spamRepository.UpdateSpamDomains()
	if err != nil {
		return UpdateResult{}, err
	}

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

// ListSpamDomains returns a list of all known spam domains
func (spamControl *SpamControl) ListSpamDomains() ([]string, error) {
	return spamControl.spamRepository.GetSpamDomains()
}

// DetectSpam checks the given account for referrer spam.
// Returns an error if the analysis failed.
func (spamControl *SpamControl) DetectSpam(accountID string, numberOfDaysToLookBack int) (AnalysisResult, error) {
	if numberOfDaysToLookBack < 1 {
		return AnalysisResult{}, fmt.Errorf("The specified number of days to look back cannot be below 1")
	}

	return spamControl.spamAnalysis.GetSpamAnalysis(accountID, numberOfDaysToLookBack, 0.75)
}

// GetTrainingData returns a set of training data for the given accountID.
// Returns an error if the training data could not be fetched.
func (spamControl *SpamControl) GetTrainingData(accountID string, numberOfDaysToLookBack int) (MachineLearningModel, error) {
	return spamControl.spamAnalysis.GetTrainingData(accountID, numberOfDaysToLookBack)
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

	knownSpamDomains, spamDomainsError := spamControl.spamRepository.GetSpamDomains()
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
			fmt.Printf("Removing filter %q\n", filterStatus.Filter().Name)
			removeError := spamControl.filterProvider.RemoveFilter(account.ID, filterStatus.Filter().ID)
			if removeError != nil {
				return removeError
			}

			continue
		}

		// update outdated filters
		if filterStatus.Status() == status.Outdated {
			fmt.Printf("Updating filter %q\n", filterStatus.Filter().Name)
			_, updateError := spamControl.filterProvider.UpdateFilter(account.ID, filterStatus.Filter().ID, filterStatus.Filter())
			if updateError != nil {
				return updateError
			}

			continue
		}

		// create new filters
		if filterStatus.Status() == status.NotInstalled {
			fmt.Printf("Creating filter %q\n", filterStatus.Filter().Name)
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

func getAverageProbability(spamDomains []SpamDomain) float64 {
	if len(spamDomains) == 0 {
		return 0.0
	}

	totalProbability := 0.0
	for _, spamDomain := range spamDomains {
		totalProbability += spamDomain.SpamProbability
	}

	numberOfDomains := float64(len(spamDomains))
	return totalProbability / numberOfDomains
}
