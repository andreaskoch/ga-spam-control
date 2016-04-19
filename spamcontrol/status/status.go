package status

import "math"

const (
	NotSet Status = 1 + iota
	Unknown
	UpToDate
	Error
	Outdated
	NotInstalled
)

var labels = []string{
	"not set",
	"unknown",
	"up-to-date",
	"error",
	"outdated",
	"not installed",
}

type Status int

func (status Status) String() string {
	return labels[status-1]
}

// CalculateGlobalStatus determines a global status
// based on the given sub-statuses.
func CalculateGlobalStatus(subStatuses []Status) Status {

	// Status: unknown
	if len(subStatuses) == 0 {
		return Unknown
	}

	// If all statuses are the same, return that.
	if statusesAreAlike, status := allStatusesAreAlike(subStatuses); statusesAreAlike {
		return status
	}

	// If there is a majority, return that.
	if hasMajority, status := getMajorityStatus(subStatuses); hasMajority {
		return status
	}

	return Unknown
}

// allStatusesAreAlike checks if all given statuses are the same.
// Returns true and the status if all statuses are alike.
// Otherwise false and nil.
func allStatusesAreAlike(statuses []Status) (bool, Status) {
	if statuses == nil || len(statuses) == 0 {
		return false, NotSet
	}

	statusUsages := make(map[Status]int)

	// build statusUsages statistic
	for _, status := range statuses {

		if value, exists := statusUsages[status]; exists {
			statusUsages[status] = value + 1
		} else {
			statusUsages[status] = 1
		}
	}

	if len(statusUsages) == 1 {
		return true, statuses[0]
	}

	return false, NotSet
}

// getMajorityNumber returns the number that
// would represent the majority for the given population.
func getMajorityNumber(population int) int {
	majority := math.Ceil(float64(population) * 0.5)
	remainder := math.Mod(float64(population), majority)
	if population != 1 && remainder == 0 {
		majority += 1
	}

	return int(majority)
}

// getMajorityStatus detects if there is a majority
// status in the given list of statuses.
// If yes, it will return true and the status that makes
// up the majority of the entries.
// If no, it will return false and nil.
func getMajorityStatus(statuses []Status) (bool, Status) {
	if statuses == nil || len(statuses) == 0 {
		return false, NotSet
	}

	statusUsages := make(map[Status]int)

	// build statusUsages statistic
	for _, status := range statuses {
		statusUsages[status] = statusUsages[status] + 1
	}

	// determine the majority
	majority := getMajorityNumber(len(statuses))

	for _, status := range statuses {
		statusUsage := statusUsages[status]
		if statusUsage >= majority {

			// majority exists
			return true, status
		}
	}

	// no majority
	return false, NotSet
}
