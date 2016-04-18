package spamcontrol

import (
	"log"
	"strings"

	"github.com/andreaskoch/ga-spam-control/api"
)

type SpamController interface {
	Remove() error
	Status() (StateOverview, error)
	Update() error
}

// New creates a new spam control instance.
func New(analyticsAPI api.AnalyticsAPI) *SpamControl {
	return &SpamControl{
		analyticsAPI: analyticsAPI,
		filterName:   "ga-spam-control",
	}
}

// The SpamControl type provides functions for
// managing Google Analtics spam filters.
type SpamControl struct {
	analyticsAPI api.AnalyticsAPI
	filterName   string
}

func (spamControl *SpamControl) Remove() error {
	// get all available accounts
	accounts, accountsError := spamControl.analyticsAPI.GetAccounts()
	if accountsError != nil {
		return accountsError
	}

	for _, account := range accounts {

		// get all filters for account
		filters, filtersError := spamControl.analyticsAPI.GetFilters(account.ID)
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
	accounts, accountsError := spamControl.analyticsAPI.GetAccounts()
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
		filters, filtersError := spamControl.analyticsAPI.GetFilters(account.ID)
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

			// ignore all non spam-control filters
			if !strings.HasPrefix(filter.Name, spamControl.filterName) {
				continue
			}

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

func (spamControl *SpamControl) Update() error {
	filter := api.Filter{
		Name: spamControl.filterName,
		Type: "EXCLUDE",
		ExcludeDetails: api.FilterDetail{
			Kind:            "analytics#filterExpression",
			Field:           "CAMPAIGN_SOURCE",
			MatchType:       "MATCHES",
			ExpressionValue: `example\.com`,
			CaseSensitive:   false,
		},
	}

	return spamControl.analyticsAPI.CreateFilter("578578", filter)
}
