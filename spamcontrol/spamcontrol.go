package spamcontrol

import (
	"fmt"
	"log"

	"github.com/andreaskoch/ga-spam-control/api"
)

type SpamController interface {
	Remove() error
	Status() (StateOverview, error)
	Update() error
}

// New creates a new spam control instance.
func New(analyticsAPI api.AnalyticsAPI) *SpamControl {

	accountProvider := remoteAccountProvider{analyticsAPI}

	domainProvider := &remoteSpamDomainProvider{}
	filterNameProvider := &spamFilterNameProvider{"ga-spam-control"}

	filterProvider := &remoteFilterProvider{
		analyticsAPI:       analyticsAPI,
		filterNameProvider: filterNameProvider,
	}

	filterFactory := &spamFilterFactory{
		filterNameProvider:   filterNameProvider,
		filterValueMaxLength: 255,
	}

	return &SpamControl{
		accountProvider: accountProvider,
		domainProvider:  domainProvider,
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
			log.Printf("%#v\n", filter)
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

	overview := &StateOverview{
		OverallStatus: StatusUnknown(),
		Accounts:      make([]AccountStatus, 0),
	}

	// get the status for each account
	overallStatusIsSet := false
	for _, account := range accounts {

		accountStatus := &AccountStatus{
			AccountID:   account.ID,
			AccountName: account.Name,
			Status:      StatusUnknown(),
		}

		// get all filters for account
		filters, filtersError := spamControl.filterProvider.GetExistingFilters(account.ID)
		if filtersError != nil {

			// failed to fetch filters for account
			// set status: error
			accountStatus.Status = StatusError(filtersError.Error())
			overview.Accounts = append(overview.Accounts, *accountStatus)

			// set overall status
			overview.OverallStatus = StatusError(filtersError.Error())
			overallStatusIsSet = true

			continue
		}

		// check if spam control filter exists
		filterContent := "example.com"
		hasSpamControlFilter := false
		spamControlFilterIsUpToDate := false

		for _, filter := range filters {

			// there is a spam control filter
			hasSpamControlFilter = true

			// check if it needs to be updated
			if filter.ExcludeDetails.ExpressionValue == filterContent {
				spamControlFilterIsUpToDate = true
			} else {
				spamControlFilterIsUpToDate = false
			}
		}

		if !hasSpamControlFilter {
			// set status: not-installed
			accountStatus.Status = StatusNotInstalled()

			if !overallStatusIsSet {
				overview.OverallStatus = StatusNotInstalled()
				overallStatusIsSet = true
			}
		}

		if isUpToDate := hasSpamControlFilter && spamControlFilterIsUpToDate; isUpToDate {
			// set status: up-to-date
			accountStatus.Status = StatusUpToDate()

			if !overallStatusIsSet {
				overview.OverallStatus = StatusUpToDate()
			}
		}

		if isOutdated := hasSpamControlFilter && !spamControlFilterIsUpToDate; isOutdated {
			// set status: outdated
			accountStatus.Status = StatusOutdated()
			overview.OverallStatus = StatusOutdated()

			if !overallStatusIsSet {
				overview.OverallStatus = StatusOutdated()
				overallStatusIsSet = true
			}
		}

		overview.Accounts = append(overview.Accounts, *accountStatus)

		// reset status
		hasSpamControlFilter = false
		spamControlFilterIsUpToDate = false

	}

	return *overview, nil
}

// Update creates or updates spam-control filters for all accounts.
func (spamControl *SpamControl) Update() error {

	// get the latest spam domain names
	spamDomainNames, spamDomainError := spamControl.domainProvider.GetSpamDomains()
	if spamDomainError != nil {
		return spamDomainError
	}

	// create new filters for the given domain names
	filters, filterError := spamControl.filterFactory.GetNewFilters(spamDomainNames)
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
			if createError := spamControl.filterProvider.CreateFilter(account.ID, filter); createError != nil {
				return fmt.Errorf("Failed to create filter for account %q (%s): %s", account.Name, account.ID, createError.Error())
			}
		}
	}

	return nil
}
