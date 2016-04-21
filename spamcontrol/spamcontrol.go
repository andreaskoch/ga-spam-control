package spamcontrol

import (
	"fmt"

	"github.com/andreaskoch/ga-spam-control/api"
	"github.com/andreaskoch/ga-spam-control/spamcontrol/status"
)

type SpamController interface {
	Remove() error
	Status() (StateOverview, error)
	Update() error
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

	return &SpamControl{
		accountProvider: accountProvider,
		filterFactory:   filterFactory,
		filterProvider:  filterProvider,
	}
}

// The SpamControl type provides functions for
// managing Google Analtics spam filters.
type SpamControl struct {
	accountProvider accountProvider
	domainProvider  spamDomainProvider
	filterFactory   filterFactory
	filterProvider  filterProvider
}

func (spamControl *SpamControl) Remove() error {
	// get all available accounts
	accounts, accountsError := spamControl.accountProvider.GetAccounts()
	if accountsError != nil {
		return accountsError
	}

	for _, account := range accounts {

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

	}

	return nil
}

// Status collects the current spam-control status of all accessible
// analytics accounts. It returns the a StateOverview model with the Status
// of all accounts and an overall Status. If the status cannot be determined
// an error will be returned.
func (spamControl *SpamControl) Status() (StateOverview, error) {
	// get all available accounts
	accounts, accountsError := spamControl.accountProvider.GetAccounts()
	if accountsError != nil {
		return StateOverview{}, accountsError
	}

	overviewModel := &StateOverview{
		OverallStatus: status.Unknown,
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

// Update creates or updates spam-control filters for all accounts.
func (spamControl *SpamControl) Update() error {

	// create new filters for the given domain names
	filters, filterError := spamControl.filterFactory.GetNewFilters()
	if filterError != nil {
		return filterError
	}

	// get all available accounts
	accounts, accountsError := spamControl.accountProvider.GetAccounts()
	if accountsError != nil {
		return accountsError
	}

	// create the filters for all accounts
	for _, account := range accounts {
		for _, filter := range filters {
			if _, createError := spamControl.filterProvider.CreateFilter(account.ID, filter); createError != nil {
				return fmt.Errorf("Failed to create filter for account %q (%s): %s", account.Name, account.ID, createError.Error())
			}
		}
	}

	return nil
}
