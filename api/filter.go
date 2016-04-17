package api

import (
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
		AccountID: source.AccountID,
		Item: apiservice.Item{
			ID:       source.ID,
			Name:     source.Name,
			Kind:     source.Kind,
			Type:     source.Type,
			SelfLink: source.Link,
		},
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

type FilterDetail struct {
	Kind            string
	Field           string
	MatchType       string
	ExpressionValue string
	CaseSensitive   bool
}
