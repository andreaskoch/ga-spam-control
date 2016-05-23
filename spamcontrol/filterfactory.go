package spamcontrol

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/andreaskoch/ga-spam-control/api"
)

// The filterFactory interface provides a function for creating
// api.Filter models from domain names.
type filterFactory interface {

	// GetNewFilters returns a list of new api.Filter models for the given domain names.
	GetNewFilters(domainNames []string) ([]api.Filter, error)
}

// The googleAnalyticsFilterFactory creates api.Filter models
// from domain names for Google Analytics.
type googleAnalyticsFilterFactory struct {
	filterNameProvider   filterNameProvider
	filterValueMaxLength int
}

// GetNewFilters returns a list of new api.Filter models for the given domain names.
func (filterFactory googleAnalyticsFilterFactory) GetNewFilters(domainNames []string) ([]api.Filter, error) {

	var filters []api.Filter

	// escape and segment the domain names
	// for the usage as the expression value.
	filterShards, segmentsError := getFilterShards(domainNames, filterFactory.filterValueMaxLength)
	if segmentsError != nil {
		return nil, segmentsError
	}

	for _, segment := range filterShards {

		for index, expresionValue := range segment.Entries {

			filter := api.Filter{
				Kind: "analytics#filter",
				Type: "EXCLUDE",
				Name: filterFactory.filterNameProvider.GetFilterName(segment.Name, index+1),
				ExcludeDetails: api.FilterDetail{
					Kind:            "analytics#filterExpression",
					Field:           "CAMPAIGN_SOURCE",
					MatchType:       "MATCHES",
					ExpressionValue: expresionValue,
					CaseSensitive:   false,
				},
			}

			filters = append(filters, filter)

		}

	}

	return filters, nil
}

// getFilterShards returns a list of filter shards from the given list of domain names;
// while considering the specified max segment size.
func getFilterShards(domainNames []string, maxSegmentSize int) ([]filterShard, error) {

	// create shards
	aToZ := regexp.MustCompile(`[a-zA-Z0-9]`)
	domainShards := make(map[string][]string)
	for _, domainName := range domainNames {
		if len(domainName) == 0 {
			continue
		}

		// determine the segment prefix/name
		shardName := ""
		prefix := domainName[0:1]
		if !aToZ.MatchString(prefix) {
			shardName = "💩"
		} else {
			shardName = strings.ToUpper(prefix)
		}

		domainShards[shardName] = append(domainShards[shardName], domainName)
	}

	// create segments
	var shards []filterShard
	for shardName, domains := range domainShards {

		// make sure the domain names are sorted
		sort.Strings(domains)

		shard := filterShard{
			Name:    shardName,
			Entries: []string{},
		}

		currentShardValue := ""
		for _, domainName := range domains {

			if !isValidDomainName(domainName) {
				return nil, fmt.Errorf("Domain names cannot be emmpty.")
			}

			currentSegmentLength := len(currentShardValue)
			escapedDomainName := regexp.QuoteMeta(domainName)
			newDomainNameLength := len(escapedDomainName)

			// check if the domain name fits into a segment
			if newDomainNameLength >= maxSegmentSize {
				return nil, fmt.Errorf("The domain name %q is too long to fit into a segment (Max length: %d).", domainName, maxSegmentSize)
			}

			// start a new shard
			if currentSegmentLength > 0 && currentSegmentLength+newDomainNameLength+1 > maxSegmentSize {
				shard.Entries = append(shard.Entries, currentShardValue)
				currentShardValue = ""
			}

			// add domain name to current segment
			if currentShardValue == "" {
				currentShardValue = escapedDomainName
			} else {
				currentShardValue = currentShardValue + "|" + escapedDomainName
			}

		}

		// add the rest
		if currentShardValue != "" {
			shard.Entries = append(shard.Entries, currentShardValue)
		}

		shards = append(shards, shard)

	}

	sortShardsBy(filterShardsByName).Sort(shards)

	return shards, nil
}

// isValidDomainName checks if the given domain name is valid or not.
func isValidDomainName(domainName string) bool {
	if domainName == "" || strings.TrimSpace(domainName) == "" {
		return false
	}

	return true
}

// A filterShard is a building block for a list of api.Filter models.
// The segment filterShard.Name is the Filter shard ('Filters for domain beginning with "A"').
// The filterShard.Entries are the filter values for each filter in that shard.
type filterShard struct {
	Name    string
	Entries []string
}

// filterShardsByName can be used to sort filterShards by name (ascending).
func filterShardsByName(filterShard1, filterShard2 filterShard) bool {
	return filterShard1.Name < filterShard2.Name
}

// The sortShardsBy function sorts ExpressionSegment objects.
type sortShardsBy func(filterShard1, filterShard2 filterShard) bool

// Sort the given ExpressionSegment objects.
func (by sortShardsBy) Sort(filterShards []filterShard) {
	sorter := &filterShardSorter{
		filterShards: filterShards,
		by:           by,
	}

	sort.Sort(sorter)
}

type filterShardSorter struct {
	filterShards []filterShard
	by           sortShardsBy
}

func (sorter *filterShardSorter) Len() int {
	return len(sorter.filterShards)
}

func (sorter *filterShardSorter) Swap(i, j int) {
	sorter.filterShards[i], sorter.filterShards[j] = sorter.filterShards[j], sorter.filterShards[i]
}

func (sorter *filterShardSorter) Less(i, j int) bool {
	return sorter.by(sorter.filterShards[i], sorter.filterShards[j])
}
