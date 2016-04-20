package spamcontrol

import (
	"github.com/andreaskoch/ga-spam-control/api"
	"github.com/andreaskoch/ga-spam-control/spamcontrol/status"
)

type filterProvider interface {
	// GetExistingFilters returns a list of all existing api.Filter models
	// for the given account ID.
	GetExistingFilters(accountID string) ([]api.Filter, error)

	// CreateFilter creates the given filter.
	CreateFilter(accountID string, filter api.Filter) (api.Filter, error)

	// RemoveFilter deletes the given filter from the specified account.
	RemoveFilter(accountID, filterID string) error

	// GetFilterStatus returns the status of the filter for the given account ID.
	GetFilterStatus(accountID string) status.Status
}

type remoteFilterProvider struct {
	analyticsAPI       api.AnalyticsAPI
	filterNameProvider filterNameProvider
	filterFactory      filterFactory
}

// GetExistingFilters returns a list of all existing api.Filter models
// for the given account ID.
func (filterProvider remoteFilterProvider) GetExistingFilters(accountID string) ([]api.Filter, error) {
	filters := make([]api.Filter, 0)

	allFilters, filtersError := filterProvider.analyticsAPI.GetFilters(accountID)
	if filtersError != nil {
		return nil, filtersError
	}

	for _, filter := range allFilters {

		// ignore all non spam-control filters
		if !filterProvider.filterNameProvider.IsValidFilterName(filter.Name) {
			continue
		}

		filters = append(filters, filter)
	}

	return filters, nil
}

// CreateFilter creates the given filter.
func (filterProvider remoteFilterProvider) CreateFilter(accountID string, filter api.Filter) (api.Filter, error) {
	return filterProvider.analyticsAPI.CreateFilter(accountID, filter)
}

// RemoveFilter deletes the given filter from the specified account.
func (filterProvider remoteFilterProvider) RemoveFilter(accountID, filterID string) error {
	return filterProvider.analyticsAPI.RemoveFilter(accountID, filterID)
}

// GetFilterStatus returns the status of the filter for the given account ID.
func (filterProvider remoteFilterProvider) GetFilterStatus(accountID string) status.Status {

	// get the existing filters
	existingFilters, existingFilterError := filterProvider.GetExistingFilters(accountID)

	// Status: error (cannot fetch existing filters)
	if existingFilterError != nil {
		return status.Error
	}

	// Status: not-installed
	if len(existingFilters) == 0 {
		return status.NotInstalled
	}

	// get the latest filters
	latestFilters, latestFiltersError := filterProvider.filterFactory.GetNewFilters()

	// Status: error (cannot determine new filters)
	if latestFiltersError != nil {
		return status.Error
	}

	// Status: outdated
	if len(existingFilters) != len(latestFilters) {
		return status.Outdated
	}

	// check content of each filter
	for index, existingFilter := range existingFilters {
		newFilter := latestFilters[index]

		// Status: outdated
		if !existingFilter.Equals(newFilter) {
			return status.Outdated
		}
	}

	return status.UpToDate
}
