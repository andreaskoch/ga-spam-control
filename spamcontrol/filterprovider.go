package spamcontrol

import (
	"github.com/andreaskoch/ga-spam-control/api"
	"github.com/andreaskoch/ga-spam-control/spamcontrol/status"
)

// A filterProvider offers CRUD operations for analytics filters.
type filterProvider interface {
	// GetExistingFilters returns a list of all existing api.Filter models
	// for the given account ID.
	GetExistingFilters(accountID string) ([]api.Filter, error)

	// CreateFilter creates the given filter.
	CreateFilter(accountID string, filter api.Filter) (api.Filter, error)

	// UpdateFilter updates the given filter.
	UpdateFilter(accountID string, filterID string, filter api.Filter) (api.Filter, error)

	// RemoveFilter deletes the given filter from the specified account.
	RemoveFilter(accountID, filterID string) error

	// GetAccountStatus returns the overall status for the given account ID.
	GetAccountStatus(accountID string) (InstallationStatus, error)

	// GetFilterStatuses returns the individual filter statuses for the given account.
	GetFilterStatuses(accountID string) (FilterStatuses, error)
}

// A remoteFilterProvider offers CRUD operations for analytics filters
// the an analytics API.
type remoteFilterProvider struct {
	analyticsAPI api.AnalyticsAPI

	spamDomainProvider SpamDomainProvider

	filterNameProvider filterNameProvider
	filterFactory      filterFactory
}

// GetExistingFilters returns a list of all existing api.Filter models
// for the given account ID.
func (filterProvider remoteFilterProvider) GetExistingFilters(accountID string) ([]api.Filter, error) {
	var filters []api.Filter

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

// UpdateFilter updates the given filter.
func (filterProvider remoteFilterProvider) UpdateFilter(accountID string, filterID string, filter api.Filter) (api.Filter, error) {
	return filterProvider.analyticsAPI.UpdateFilter(accountID, filterID, filter)
}

// RemoveFilter deletes the given filter from the specified account.
func (filterProvider remoteFilterProvider) RemoveFilter(accountID, filterID string) error {
	return filterProvider.analyticsAPI.RemoveFilter(accountID, filterID)
}

// GetAccountStatus returns the status of for the given account ID.
func (filterProvider remoteFilterProvider) GetAccountStatus(accountID string) (InstallationStatus, error) {

	filterStatuses, filterStatusError := filterProvider.GetFilterStatuses(accountID)
	if filterStatusError != nil {
		return InstallationStatus{}, filterStatusError
	}

	return filterStatuses.OverallStatus(), nil
}

// GetFilterStatuses returns the individual filter statuses for the given account.
func (filterProvider remoteFilterProvider) GetFilterStatuses(accountID string) (FilterStatuses, error) {
	// get the existing filters
	existingFilters, existingFilterError := filterProvider.GetExistingFilters(accountID)
	if existingFilterError != nil {
		return nil, existingFilterError
	}

	// get the latest referrer spam domain names
	domainNames, spamDomainProviderError := filterProvider.spamDomainProvider.GetSpamDomains()
	if spamDomainProviderError != nil {
		return nil, spamDomainProviderError
	}

	// get the latest filters
	latestFilters, latestFiltersError := filterProvider.filterFactory.GetNewFilters(domainNames)
	if latestFiltersError != nil {
		return nil, latestFiltersError
	}

	filterStatuses := getFilterStatuses(existingFilters, latestFilters)
	return filterStatuses, nil
}

// getFilterStatuses returns an overview of the Status of all given filters.
func getFilterStatuses(existingFilters, latestFilters []api.Filter) FilterStatuses {

	statuses := make(FilterStatuses, 0)

	// create an index
	oldFilters := getFilterNameMap(existingFilters)
	newFilters := getFilterNameMap(latestFilters)

	for oldName, oldFilter := range oldFilters {

		if newFilter, filterStillExists := newFilters[oldName]; filterStillExists {

			if oldFilter.Equals(newFilter) {
				// Status: up-to-date
				statuses = append(statuses, newFilterStatus(oldFilter, status.UpToDate))
				continue
			}

			// Status: outdated
			updateFilter := newFilter
			updateFilter.ID = oldFilter.ID
			statuses = append(statuses, newFilterStatus(updateFilter, status.Outdated))
			continue
		}

		// Status: obsolete
		statuses = append(statuses, newFilterStatus(oldFilter, status.Obsolete))

	}

	for newName, newFilter := range newFilters {
		if _, filterStillExists := oldFilters[newName]; filterStillExists {
			// if the filter still exists we catched it in the previous round
			continue
		}

		// Status: not-installed
		statuses = append(statuses, newFilterStatus(newFilter, status.NotInstalled))
	}

	// sort
	SortFiltersBy(filterStatusesByName).Sort(statuses)

	return statuses
}

// getFilterNameMap groups a slice api.Filter models by their name.
func getFilterNameMap(filters []api.Filter) map[string]api.Filter {
	nameMap := make(map[string]api.Filter)
	for _, filter := range filters {
		nameMap[filter.Name] = filter
	}

	return nameMap
}
