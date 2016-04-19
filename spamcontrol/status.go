package spamcontrol

import "math"

// A Status defines the current situation of the spam-control filters.
// A Status has always has a name and potentially a description.
type Status interface {
	// String returns the name of the Status.
	String() string

	// Name returns the name of the Status.
	Name() string

	// Details returns the description text of the Status (optional).
	Details() string

	// Equals check if the current status matches to given Status.
	Equals(otherStatus Status) bool
}

type baseStatus struct {
	name    string
	details string
}

// String returns a string-representation of the Status.
func (status baseStatus) String() string {
	return status.name
}

// Name returns the name of the Status.
func (status baseStatus) Name() string {
	return status.name
}

// Details returns the description text of the Status.
func (status baseStatus) Details() string {
	return status.details
}

// Equals check if the current status matches to given Status.
func (status baseStatus) Equals(otherStatus Status) bool {
	return status.Name() == otherStatus.Name()
}

type Unknown struct {
	baseStatus
}

type UpToDate struct {
	baseStatus
}

type Error struct {
	baseStatus
}

type Outdated struct {
	baseStatus
}

type NotInstalled struct {
	baseStatus
}

// StatusUnknown creates a new "unknown" Status instance.
// This Status type can be used as the default Status.
func StatusUnknown() Status {
	return Unknown{baseStatus: baseStatus{"unknown", ""}}
}

// StatusUpToDate creates a new "up-to-date" Status instance.
// This Status type can be used when all spam-control mechanisms
// are installed in the latest available version.
func StatusUpToDate() Status {
	return UpToDate{baseStatus: baseStatus{"up-to-date", ""}}
}

// StatusError creates a new "error" Status instance.
// The given errorMessage will be assigned to the Status.details.
// This Status type can be used if an error occurred while
// determining the status.
func StatusError(errorMessage string) Status {
	return Error{baseStatus: baseStatus{"error", errorMessage}}
}

// StatusOutdated creates a new "outdated" Status instance.
// This Status type can be used when all spam-control mechanisms
// are installed but not in the latest available version.
func StatusOutdated() Status {
	return Outdated{baseStatus: baseStatus{"outdated", ""}}
}

// StatusOutdated creates a new "not-installed" Status instance.
// This Status type can be used when no spam-control mechanisms
// are installed.
func StatusNotInstalled() Status {
	return NotInstalled{baseStatus: baseStatus{"not-installed", ""}}
}

// calculateGlobalStatus determines a global status
// based on the given sub-statuses.
func calculateGlobalStatus(subStatuses []Status) Status {

	// Status: unknown
	if len(subStatuses) == 0 {
		return StatusUnknown()
	}

	// If all statuses are the same, return that.
	if statusesAreAlike, status := allStatusesAreAlike(subStatuses); statusesAreAlike {
		return status
	}

	// If there is a majority, return that.
	if hasMajority, status := getMajorityStatus(subStatuses); hasMajority {
		return status
	}

	return StatusUnknown()
}

// allStatusesAreAlike checks if all given statuses are the same.
// Returns true and the status if all statuses are alike.
// Otherwise false and nil.
func allStatusesAreAlike(statuses []Status) (bool, Status) {
	if statuses == nil || len(statuses) == 0 {
		return false, nil
	}

	statusUsages := make(map[string]int)

	// build statusUsages statistic
	for _, status := range statuses {
		if status == nil {
			continue
		}

		if value, exists := statusUsages[status.Name()]; exists {
			statusUsages[status.Name()] = value + 1
		} else {
			statusUsages[status.Name()] = 1
		}
	}

	if len(statusUsages) == 1 {
		return true, statuses[0]
	}

	return false, nil
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
		return false, nil
	}

	statusUsages := make(map[string]int)

	// build statusUsages statistic
	for _, status := range statuses {
		statusUsages[status.Name()] = statusUsages[status.Name()] + 1
	}

	// determine the majority
	majority := getMajorityNumber(len(statuses))

	for _, status := range statuses {
		statusUsage := statusUsages[status.Name()]
		if statusUsage >= majority {

			// majority exists
			return true, status
		}
	}

	// no majority
	return false, nil
}
