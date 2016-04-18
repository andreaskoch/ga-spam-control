package spamcontrol

import (
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
	return &SpamControl{
		analyticsAPI: analyticsAPI,
	}
}

// The SpamControl type provides functions for
// managing Google Analtics spam filters.
type SpamControl struct {
	analyticsAPI api.AnalyticsAPI
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
		OverallStatus: StatusUpToDate(),
	}

	for _, account := range accounts {

		// get all filters for account
		_, filtersError := spamControl.analyticsAPI.GetFilters(account.ID)
		if filtersError != nil {
			return StateOverview{}, filtersError
		}

		// for _, filter := range filters {
		// log.Printf("%#v\n", filter)
		// }

	}

	return *overview, nil
}

func (spamControl *SpamControl) Update() error {
	filter := api.Filter{
		Name: "jkljk",
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
