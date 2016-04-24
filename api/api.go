package api

import (
	"fmt"

	"github.com/andreaskoch/ga-spam-control/api/apicredentials"
	"github.com/andreaskoch/ga-spam-control/api/apiservice"
)

// The AnalyticsAPI interface provides Analytics API functions.
type AnalyticsAPI interface {

	// GetAccounts returns all apiservice.Account models.
	GetAccounts() ([]Account, error)

	// GetAnalyticsData analytics data for the given account ID.
	GetAnalyticsData(accountID string) (AnalyticsData, error)

	// CreateFilter creates a new Filter for the given account ID.
	CreateFilter(accountID string, filter Filter) (Filter, error)

	// UpdateFilter updates the given filter.
	UpdateFilter(accountID string, filterID string, filter Filter) (Filter, error)

	// GetFilters returns all Filter models for the given account.
	GetFilters(accountID string) ([]Filter, error)

	// GetProfiles returns all Profile models for the given account.
	GetProfiles(accountID string) ([]Profile, error)

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

// GetAnalyticsData analytics data for the given account ID.
func (api *API) GetAnalyticsData(accountID string) (AnalyticsData, error) {

	// get the profiles to which the filter will be assigned
	profiles, profilesError := api.GetProfiles(accountID)
	if profilesError != nil {
		return AnalyticsData{}, profilesError
	}

	if len(profiles) == 0 {
		return AnalyticsData{}, fmt.Errorf("No profiles found for account %q", accountID)
	}

	profile := profiles[0]
	serviceData, analyticsDataErr := api.service.GetAnalyticsData(profile.ID)
	if analyticsDataErr != nil {
		return AnalyticsData{}, analyticsDataErr
	}

	return toModelAnalyticsData(serviceData), nil
}

// CreateFilter creates a new Filter for the given account ID.
func (api *API) CreateFilter(accountID string, filterParameters Filter) (Filter, error) {

	// get the profiles to which the filter will be assigned
	profiles, profilesError := api.GetProfiles(accountID)
	if profilesError != nil {
		return Filter{}, profilesError
	}

	// create the filter
	serviceFilter, err := api.service.CreateFilter(accountID, toServiceFilter(filterParameters))
	if err != nil {
		return Filter{}, err
	}

	// link the created filter to all available profiles
	for _, profile := range profiles {
		err := api.service.CreateProfileFilterLink(accountID, profile.ID, profile.WebPropertyID, serviceFilter.ID)
		if err != nil {
			return Filter{}, err
		}
	}

	createdFilter := toModelFilter(serviceFilter)
	return createdFilter, nil
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

// GetProfiles returns all Profile models for the given account.
func (api *API) GetProfiles(accountID string) ([]Profile, error) {
	serviceProfiles, err := api.service.GetProfiles(accountID)
	if err != nil {
		return nil, err
	}

	profiles := toModelProfiles(serviceProfiles)

	return profiles, nil
}

// UpdateFilter updates the given filter.
func (api *API) UpdateFilter(accountID string, filterID string, filterParameters Filter) (Filter, error) {

	serviceFilter, err := api.service.UpdateFilter(accountID, filterID, toServiceFilter(filterParameters))
	if err != nil {
		return Filter{}, err
	}

	return toModelFilter(serviceFilter), nil
}

// RemoveFilter deletes the given filter from the specified account.
func (api *API) RemoveFilter(accountID, filterID string) error {

	err := api.service.RemoveFilter(accountID, filterID)
	if err != nil {
		return err
	}

	return nil
}
