package api

import (
	"github.com/andreaskoch/ga-spam-control/api/apicredentials"
	"github.com/andreaskoch/ga-spam-control/api/apiservice"
)

// The AnalyticsAPI interface provides Analytics API functions.
type AnalyticsAPI interface {

	// GetAccounts returns all apiservice.Account models.
	GetAccounts() ([]Account, error)

	// CreateFilter creates a new Filter for the given account ID.
	CreateFilter(accountID string, filter Filter) error

	// GetFilters returns all Filter models for the given account.
	GetFilters(accountID string) ([]Filter, error)

	// RemoveFilter deletes the given filter from the specified account.
	RemoveFilter(accountID, filterID string) error
}

// New creates a new API instance.
func New(tokenStore apicredentials.TokenStorer, clientID, clientSecret string) (AnalyticsAPI, error) {
	service, serviceError := apiservice.New(tokenStore, clientID, clientSecret)
	if serviceError != nil {
		return nil, serviceError
	}

	return &API{
		service: service,
	}, nil
}

// API provides CRUD operations for the Google Analytics API.
type API struct {
	service *apiservice.GoogleAnalytics
}

// GetAccounts returns all apiservice.Account models.
func (api *API) GetAccounts() ([]Account, error) {
	serviceAccounts, err := api.service.GetAccounts()
	if err != nil {
		return nil, err
	}

	accounts := toModelAccounts(serviceAccounts)

	SortAccountsBy(accountsByID).Sort(accounts)

	return accounts, nil
}

// CreateFilter creates a new Filter for the given account ID.
func (api *API) CreateFilter(accountID string, filter Filter) error {

	err := api.service.CreateFilter(accountID, toServiceFilter(filter))
	if err != nil {
		return err
	}

	return nil
}

// GetFilters returns all Filter models for the given account.
func (api *API) GetFilters(accountID string) ([]Filter, error) {
	serviceFilters, err := api.service.GetFilters(accountID)
	if err != nil {
		return nil, err
	}

	filters := toModelFilters(serviceFilters)

	SortFiltersBy(filtersByName).Sort(filters)

	return filters, nil
}

// RemoveFilter deletes the given filter from the specified account.
func (api *API) RemoveFilter(accountID, filterID string) error {

	err := api.service.RemoveFilter(accountID, filterID)
	if err != nil {
		return err
	}

	return nil
}
