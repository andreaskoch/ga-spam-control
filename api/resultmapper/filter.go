package resultmapper

import (
	"github.com/andreaskoch/ga-spam-control/api/apimodel"
	"github.com/andreaskoch/ga-spam-control/api/apiservice"
)

func MapFilters(sources []apiservice.Filter) []apimodel.Filter {

	accounts := make([]apimodel.Filter, 0)
	for _, source := range sources {
		accounts = append(accounts, MapFilter(source))
	}

	return accounts
}

// MapFilter converts a apiservice.Filter model into a apimodel.Filter model.
func MapFilter(source apiservice.Filter) apimodel.Filter {
	return apimodel.Filter{
		AccountID: source.AccountID,
		ID:        source.ID,
		Name:      source.Name,
		Kind:      source.Kind,
		Type:      source.Type,
		Link:      source.SelfLink,
		ExcludeDetails: apimodel.FilterDetail{
			Kind:            source.ExcludeDetails.Kind,
			Field:           source.ExcludeDetails.Field,
			MatchType:       source.ExcludeDetails.MatchType,
			ExpressionValue: source.ExcludeDetails.ExpressionValue,
			CaseSensitive:   source.ExcludeDetails.CaseSensitive,
		},
	}
}
