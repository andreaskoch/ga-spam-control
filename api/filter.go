package api

import (
	"sort"

	"github.com/andreaskoch/ga-spam-control/api/apiservice"
)

// toModelFilters converts []apiservice.Filter to []Filter.
func toModelFilters(sources []apiservice.Filter) []Filter {

	accounts := make([]Filter, 0)
	for _, source := range sources {
		accounts = append(accounts, toModelFilter(source))
	}

	return accounts
}

// toModelFilter converts a apiservice.Filter model into a Filter model.
func toModelFilter(source apiservice.Filter) Filter {
	return Filter{
		AccountID: source.AccountID,
		ID:        source.ID,
		Name:      source.Name,
		Kind:      source.Kind,
		Type:      source.Type,
		Link:      source.SelfLink,
		ExcludeDetails: FilterDetail{
			Kind:            source.ExcludeDetails.Kind,
			Field:           source.ExcludeDetails.Field,
			MatchType:       source.ExcludeDetails.MatchType,
			ExpressionValue: source.ExcludeDetails.ExpressionValue,
			CaseSensitive:   source.ExcludeDetails.CaseSensitive,
		},
	}
}

// toServiceFilter converts Filter to apiservice.Filter.
func toServiceFilter(source Filter) apiservice.Filter {
	return apiservice.Filter{
		Item: apiservice.Item{
			ID:       source.ID,
			Kind:     source.Kind,
			SelfLink: source.Link,
		},
		AccountID: source.AccountID,
		Name:      source.Name,
		Type:      source.Type,
		ExcludeDetails: apiservice.FilterDetail{
			Kind:            source.ExcludeDetails.Kind,
			Field:           source.ExcludeDetails.Field,
			MatchType:       source.ExcludeDetails.MatchType,
			ExpressionValue: source.ExcludeDetails.ExpressionValue,
			CaseSensitive:   source.ExcludeDetails.CaseSensitive,
		},
	}
}

type Filter struct {
	AccountID      string
	ID             string
	Kind           string
	Name           string
	Type           string
	Link           string
	ExcludeDetails FilterDetail
}

// Equals checks the this Filter matches the given one.
func (filter Filter) Equals(other Filter) bool {
	if filter.Kind != other.Kind {
		return false
	}

	if filter.Type != other.Type {
		return false
	}

	if filter.Name != other.Name {
		return false
	}

	if !filter.ExcludeDetails.Equals(other.ExcludeDetails) {
		return false
	}

	return true
}

type FilterDetail struct {
	Kind            string
	Field           string
	MatchType       string
	ExpressionValue string
	CaseSensitive   bool
}

// Equals checks the this FilterDetail matches the given one.
func (details FilterDetail) Equals(other FilterDetail) bool {
	if details.Kind != other.Kind {
		return false
	}

	if details.Field != other.Field {
		return false
	}

	if details.MatchType != other.MatchType {
		return false
	}

	if details.ExpressionValue != other.ExpressionValue {
		return false
	}

	if details.CaseSensitive != other.CaseSensitive {
		return false
	}

	return true
}

// filtersByName can be used to sort filters by name (ascending).
func filtersByName(filter1, filter2 Filter) bool {
	return filter1.Name < filter2.Name
}

type SortFiltersBy func(filter1, filter2 Filter) bool

func (by SortFiltersBy) Sort(filters []Filter) {
	sorter := &filterSorter{
		filters: filters,
		by:      by,
	}

	sort.Sort(sorter)
}

type filterSorter struct {
	filters []Filter
	by      SortFiltersBy
}

func (sorter *filterSorter) Len() int {
	return len(sorter.filters)
}

func (sorter *filterSorter) Swap(i, j int) {
	sorter.filters[i], sorter.filters[j] = sorter.filters[j], sorter.filters[i]
}

func (sorter *filterSorter) Less(i, j int) bool {
	return sorter.by(sorter.filters[i], sorter.filters[j])
}
