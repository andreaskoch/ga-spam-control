package api

import (
	"github.com/andreaskoch/ga-spam-control/api/apicredentials"
	"github.com/andreaskoch/ga-spam-control/api/apimodel"
	"github.com/andreaskoch/ga-spam-control/api/apiservice"
	"github.com/andreaskoch/ga-spam-control/api/resultmapper"
)

func New(tokenStore apicredentials.TokenStorer, clientID, clientSecret string) (*API, error) {
	service, serviceError := apiservice.New(tokenStore, clientID, clientSecret)
	if serviceError != nil {
		return nil, serviceError
	}

	return &API{
		service: service,
	}, nil
}

type API struct {
	service *apiservice.GoogleAnalytics
}

func (api *API) GetAccounts() ([]apimodel.Account, error) {
	serviceAccounts, err := api.service.GetAccounts()
	if err != nil {
		return nil, err
	}

	apiAccounts := resultmapper.MapAccounts(serviceAccounts)
	return apiAccounts, nil
}

func (api *API) GetFilters(accountID string) ([]apimodel.Filter, error) {
	serviceFilters, err := api.service.GetFilters(accountID)
	if err != nil {
		return nil, err
	}

	apiFilters := resultmapper.MapFilters(serviceFilters)
	return apiFilters, nil
}

// CreateFilter creates a new filter for the given account ID.
func (api *API) CreateFilter(accountID string) error {

	filter := apiservice.Filter{}
	filter.Name = "jkljk"
	filter.Type = "EXCLUDE"
	filter.ExcludeDetails = apiservice.FilterDetail{
		Kind:            "analytics#filterExpression",
		Field:           "CAMPAIGN_SOURCE",
		MatchType:       "MATCHES",
		ExpressionValue: `example\.com`,
		CaseSensitive:   false,
	}

	err := api.service.CreateFilter(accountID, filter)
	if err != nil {
		return err
	}

	return nil
}
