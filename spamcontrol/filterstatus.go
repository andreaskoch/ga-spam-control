package spamcontrol

import (
	"fmt"
	"sort"

	"github.com/andreaskoch/ga-spam-control/api"
	"github.com/andreaskoch/ga-spam-control/spamcontrol/status"
)

// FilterStatuses are a list of of FilterStatus objects.
type FilterStatuses []FilterStatus

// OverallStatus calculates an overall status for all status in this list.
func (filterStatuses FilterStatuses) OverallStatus() InstallationStatus {
	var numberOfFiltersWithStateUpToDate int
	for _, filterStatus := range filterStatuses {
		if filterStatus.Status() == status.UpToDate {
			numberOfFiltersWithStateUpToDate++
		}
	}

	return InstallationStatus{
		TotalFilters:    len(filterStatuses),
		UpToDateFilters: numberOfFiltersWithStateUpToDate,
	}
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
