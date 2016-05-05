package spamcontrol

import (
	"fmt"
	"sort"

	"github.com/andreaskoch/ga-spam-control/api"
	"github.com/andreaskoch/ga-spam-control/spamcontrol/status"
)

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
	GetAccountStatus(accountID string) (status.Status, error)

	// GetFilterStatuses returns the individual filter statuses for the given account.
	GetFilterStatuses(accountID string) (FilterStatuses, error)
}

type remoteFilterProvider struct {
	analyticsAPI       api.AnalyticsAPI
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
func (filterProvider remoteFilterProvider) GetAccountStatus(accountID string) (status.Status, error) {

	filterStatuses, filterStatusError := filterProvider.GetFilterStatuses(accountID)
	if filterStatusError != nil {
		return status.NotSet, filterStatusError
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

	// get the latest filters
	latestFilters, latestFiltersError := filterProvider.filterFactory.GetNewFilters()
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

func getFilterNameMap(filters []api.Filter) map[string]api.Filter {
	nameMap := make(map[string]api.Filter)
	for _, filter := range filters {
		nameMap[filter.Name] = filter
	}

	return nameMap
}

// FilterStatuses are a list of of FilterStatus objects.
type FilterStatuses []FilterStatus

// OverallStatus calculates an overall status for all status in this list.
func (filterStatuses FilterStatuses) OverallStatus() status.Status {
	var statuses []status.Status
	for _, filterStatus := range filterStatuses {
		statuses = append(statuses, filterStatus.Status())
	}

	return status.CalculateGlobalStatus(statuses)
}

// newFilterStatus creates a new FilterStatus instance.
func newFilterStatus(filter api.Filter, status status.Status) FilterStatus {
	return FilterStatus{filter, status}
}

// FilterStatus represents the status for a given api.Filter.
type FilterStatus struct {
	filter api.Filter
	status status.Status
}

func (filterStatus FilterStatus) String() string {
	return fmt.Sprintf("%s (%s)", filterStatus.filter.ID, filterStatus.status)
}

// Filter returns the api.Filter.
func (filterStatus FilterStatus) Filter() api.Filter {
	return filterStatus.filter
}

// Status returns the status.Status.
func (filterStatus FilterStatus) Status() status.Status {
	return filterStatus.status
}

// filterStatusesByName can be used to sort filterStatuses by name (ascending).
func filterStatusesByName(filterStatus1, filterStatus2 FilterStatus) bool {
	return filterStatus1.Filter().Name < filterStatus2.Filter().Name
}

// The SortFiltersBy function sorts two FilterStatus objects.
type SortFiltersBy func(filter1, filter2 FilterStatus) bool

// Sort a list of FilterStatus objects.
func (by SortFiltersBy) Sort(filterStatuses []FilterStatus) {
	sorter := &filterStatusSorter{
		filterStatuses: filterStatuses,
		by:             by,
	}

	sort.Sort(sorter)
}

type filterStatusSorter struct {
	filterStatuses []FilterStatus
	by             SortFiltersBy
}

func (sorter *filterStatusSorter) Len() int {
	return len(sorter.filterStatuses)
}

func (sorter *filterStatusSorter) Swap(i, j int) {
	sorter.filterStatuses[i], sorter.filterStatuses[j] = sorter.filterStatuses[j], sorter.filterStatuses[i]
}

func (sorter *filterStatusSorter) Less(i, j int) bool {
	return sorter.by(sorter.filterStatuses[i], sorter.filterStatuses[j])
}
