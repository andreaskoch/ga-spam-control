package spamcontrol

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

	// Status: up-to-date
	if yes, _ := allStatusesAre(subStatuses, StatusUpToDate()); yes {
		return StatusUpToDate()
	}

	// Status: outdated
	if yes, _ := allStatusesAre(subStatuses, StatusOutdated()); yes {
		return StatusOutdated()
	}

	// Status: not-installed
	if yes, _ := allStatusesAre(subStatuses, StatusNotInstalled()); yes {
		return StatusNotInstalled()
	}

	return StatusError("")
}

func allStatusesAre(statuses []Status, status Status) (yes bool, deviantStatus Status) {
	for _, subStatus := range statuses {
		if !subStatus.Equals(status) {
			return false, subStatus
		}
	}

	return true, nil
}
