package spamcontrol

import (
	"fmt"
	"strings"
)

// The filterNameProvider interface provides
// a function for creating filter names.
type filterNameProvider interface {
	// IsValidFilterName returns true if the given filterName
	// is valid or not.
	IsValidFilterName(filerName string) bool

	// GetFilterName return a filter name for the given segment name and entry index.
	GetFilterName(segmentName string, entryIndex int) string
}

// A spamFilterNameProvider provides
// spam filter names and decides if a given
// filter name is a valid name for a filter.
type spamFilterNameProvider struct {
	filterNamePrefix string
}

// IsValidFilterName returns true if the given filterName
// belongs to a spam filter.
func (nameProvider spamFilterNameProvider) IsValidFilterName(filerName string) bool {
	return strings.HasPrefix(filerName, nameProvider.filterNamePrefix)
}

// GetFilterName return a filter name for the given segment name and entry index.
func (nameProvider spamFilterNameProvider) GetFilterName(segmentName string, entryIndex int) string {
	return fmt.Sprintf("%s Segment %s #%03d", nameProvider.filterNamePrefix, segmentName, entryIndex)
}
