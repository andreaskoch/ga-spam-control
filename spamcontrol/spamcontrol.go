package spamcontrol

import (
	"fmt"

	"github.com/andreaskoch/ga-spam-control/api"
	"github.com/andreaskoch/ga-spam-control/spamcontrol/detector"
	"github.com/andreaskoch/ga-spam-control/spamcontrol/status"
)

// The SpamController interface provides functions for displaying, updating and removing referrer spam controls.
type SpamController interface {

	// Remove the referrer spam controls from the account with the given accountID.
	// Returns an error if the removal failed.
	Remove(accountID string) error

	// Analyze checks the given account for referrer spam.
	// Returns an error if the analysis failed.
	Analyze(accountID string) error

	// Status collects the current spam-control status of all accessible
	// analytics accounts. It returns the a StateOverview model with the Status
	// of all accounts and an overall Status. If the status cannot be determined
	// an error will be returned.
	GlobalStatus() (StateOverview, error)

	// AccountStatus returns the current spam-control status of the account
	// with the given account ID. Returns an error if the status cannot be
	// determined.
	AccountStatus(accountID string) (status.Status, error)

	// Update the referrer spam controls for the account with the given accountID.
	// Returns an error if the update failed.
	Update(accountID string) error
}

// New creates a new spam control instance.
func New(analyticsAPI api.AnalyticsAPI) *SpamControl {

	accountProvider := remoteAccountProvider{analyticsAPI}

	domainProvider := &remoteSpamDomainProvider{"https://raw.githubusercontent.com/ddofborg/analytics-ghost-spam-list/master/adwordsrobot.com-spam-list.txt"}
	filterNameProvider := &spamFilterNameProvider{"ga-spam-control"}

	filterFactory := &spamFilterFactory{
		domainProvider:       domainProvider,
		filterNameProvider:   filterNameProvider,
		filterValueMaxLength: 255,
	}

	filterProvider := &remoteFilterProvider{
		analyticsAPI:       analyticsAPI,
		filterNameProvider: filterNameProvider,
		filterFactory:      filterFactory,
	}

	analyticsDataProvider := &remoteAnalyticsDataProvider{
		analyticsAPI: analyticsAPI,
	}

	spamDetector := &detector.AzureMLSpamDetection{}

	return &SpamControl{
		accountProvider:       accountProvider,
		filterFactory:         filterFactory,
		filterProvider:        filterProvider,
		analyticsDataProvider: analyticsDataProvider,
		spamDetector:          spamDetector,
	}
}

// The SpamControl type provides functions for
// managing Google Analtics spam filters.
type SpamControl struct {
	accountProvider       accountProvider
	domainProvider        spamDomainProvider
	filterFactory         filterFactory
	filterProvider        filterProvider
	analyticsDataProvider analyticsDataProvider
	spamDetector          detector.SpamDetector
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

// Analyze checks the given account for referrer spam.
// Returns an error if the analysis failed.
func (spamControl *SpamControl) Analyze(accountID string) error {

	analyticsData, analyticsDataError := spamControl.analyticsDataProvider.GetAnalyticsData(accountID)
	if analyticsDataError != nil {
		return analyticsDataError
	}

	ratedAnalyticsData, spamDetectionError := spamControl.spamDetector.GetSpamRating(analyticsData)
	if spamDetectionError != nil {
		return spamDetectionError
	}

	// get all spam domains
	spamDomains := make(map[string]int)
	for _, row := range ratedAnalyticsData {
		if !row.IsSpam {
			continue
		}

		spamDomains[row.Source]++
	}

	// print a spam domain names
	for spamDomain, _ := range spamDomains {
		fmt.Println(spamDomain)
	}

	return nil
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

	overviewModel := &StateOverview{
		OverallStatus: status.NotSet,
		Accounts:      make([]AccountStatus, 0),
	}

	// get the status for each account
	accountStatuses := make([]status.Status, 0, len(accounts))
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

		// capture the account status for the calculation of the
		// overall status
		accountStatuses = append(accountStatuses, accountStatus)

		overviewModel.Accounts = append(overviewModel.Accounts, accountStatusModel)
	}

	// set the overall status
	overviewModel.OverallStatus = status.CalculateGlobalStatus(accountStatuses)

	return *overviewModel, nil
}

// AccountStatus returns the current spam-control status of the account
// with the given account ID. Returns an error if the status cannot be
// determined.
func (spamControl *SpamControl) AccountStatus(accountID string) (status.Status, error) {
	// get the requested account
	account, accountError := spamControl.accountProvider.GetAccount(accountID)
	if accountError != nil {
		return status.NotSet, accountError
	}

	// get the accounts' status
	accountStatus, accountStatusError := spamControl.filterProvider.GetAccountStatus(account.ID)
	if accountStatusError != nil {
		return status.NotSet, accountStatusError
	}

	return accountStatus, nil
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
			removeError := spamControl.filterProvider.RemoveFilter(account.ID, filterStatus.Filter().ID)
			if removeError != nil {
				return removeError
			}

			continue
		}

		// update outdated filters
		if filterStatus.Status() == status.Outdated {
			_, updateError := spamControl.filterProvider.UpdateFilter(account.ID, filterStatus.Filter().ID, filterStatus.Filter())
			if updateError != nil {
				return updateError
			}

			continue
		}

		// create new filters
		if filterStatus.Status() == status.NotInstalled {
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
