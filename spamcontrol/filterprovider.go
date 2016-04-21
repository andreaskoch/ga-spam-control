package spamcontrol

import (
	"fmt"

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

	// GetAccountStatus returns the status of the filter for the given account ID.
	GetAccountStatus(accountID string) (status.Status, error)
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

// GetAccountStatus returns the status of for the given account ID.
func (filterProvider remoteFilterProvider) GetAccountStatus(accountID string) (status.Status, error) {

	// get the existing filters
	existingFilters, existingFilterError := filterProvider.GetExistingFilters(accountID)
	if existingFilterError != nil {
		return status.Unknown, existingFilterError
	}

	// get the latest filters
	latestFilters, latestFiltersError := filterProvider.filterFactory.GetNewFilters()
	if latestFiltersError != nil {
		return status.Unknown, latestFiltersError
	}

	filterStatuses := getFilterStatuses(existingFilters, latestFilters)
	return filterStatuses.OverallStatus(), nil
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
			statuses = append(statuses, newFilterStatus(oldFilter, status.Outdated))
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

	return statuses
}

func getFilterNameMap(filters []api.Filter) map[string]api.Filter {
	nameMap := make(map[string]api.Filter)
	for _, filter := range filters {
		nameMap[filter.Name] = filter
	}

	return nameMap
}

type FilterStatuses []FilterStatus

func (filterStatuses FilterStatuses) OverallStatus() status.Status {
	statuses := make([]status.Status, 0)
	for _, filterStatus := range filterStatuses {
		statuses = append(statuses, filterStatus.Status())
	}

	return status.CalculateGlobalStatus(statuses)
}

// newFilterStatus creates a new FilterStatus instance.
func newFilterStatus(filter api.Filter, status status.Status) FilterStatus {
	return FilterStatus{filter, status}
}

type FilterStatus struct {
	filter api.Filter
	status status.Status
}

func (filterStatus FilterStatus) String() string {
	return fmt.Sprintf("%s (%s)", filterStatus.filter.ID, filterStatus.status)
}

func (filterStatus FilterStatus) Filter() api.Filter {
	return filterStatus.filter
}

func (filterStatus FilterStatus) Status() status.Status {
	return filterStatus.status
}
