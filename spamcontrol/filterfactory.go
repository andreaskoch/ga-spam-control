package spamcontrol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andreaskoch/ga-spam-control/api"
)

type filterFactory interface {

	// GetNewFilters returns a list of new api.Filter models
	// for the given list of spam domain names.
	GetNewFilters(domainNames []string) ([]api.Filter, error)
}

type spamFilterFactory struct {
	filterNameProvider   filterNameProvider
	filterValueMaxLength int
}

// GetNewFilters returns a list of new api.Filter models
// for the given list of spam domain names.
func (filterFactory spamFilterFactory) GetNewFilters(domainNames []string) ([]api.Filter, error) {

	filters := make([]api.Filter, 0)

	// escape and segment the domain names
	// for the usage as the expression value.
	expressionSegments, segmentsError := getExpressionSegments(domainNames, filterFactory.filterValueMaxLength)
	if segmentsError != nil {
		return nil, segmentsError
	}

	for index, expressionSegment := range expressionSegments {

		filter := api.Filter{
			Name: filterFactory.filterNameProvider.GetFilterName(index + 1),
			Type: "EXCLUDE",
			ExcludeDetails: api.FilterDetail{
				Kind:            "analytics#filterExpression",
				Field:           "CAMPAIGN_SOURCE",
				MatchType:       "MATCHES",
				ExpressionValue: expressionSegment,
				CaseSensitive:   false,
			},
		}

		filters = append(filters, filter)

	}

	return filters, nil
}

// getExpressionSegments returns a list of regular expression segments
// from the given list of domain names; respecting the specified
// max segment size.
func getExpressionSegments(domainNames []string, maxSegmentSize int) ([]string, error) {

	var valueSegments []string
	currentSegment := ""
	for _, domainName := range domainNames {

		if !isValidDomainName(domainName) {
			return nil, fmt.Errorf("Domain names cannot be emmpty.")
		}

		currentSegmentLength := len(currentSegment)
		escapedDomainName := regexp.QuoteMeta(domainName)
		newDomainNameLength := len(escapedDomainName)

		// check if the domain name fits into a segment
		if newDomainNameLength >= maxSegmentSize {
			return nil, fmt.Errorf("The domain name %q is too long to fit into a segment (Max length: %d).", domainName, maxSegmentSize)
		}

		// start a new segment
		if currentSegmentLength > 0 && currentSegmentLength+newDomainNameLength+1 > maxSegmentSize {
			valueSegments = append(valueSegments, currentSegment)
			currentSegment = ""
		}

		// add domain name to current segment
		if currentSegment == "" {
			currentSegment = escapedDomainName
		} else {
			currentSegment = currentSegment + "|" + escapedDomainName
		}

	}

	// add the rest
	if currentSegment != "" {
		valueSegments = append(valueSegments, currentSegment)
	}

	return valueSegments, nil
}

// isValidDomainName checks if the given domain name is valid or not.
func isValidDomainName(domainName string) bool {
	if domainName == "" || strings.TrimSpace(domainName) == "" {
		return false
	}

	return true
}
