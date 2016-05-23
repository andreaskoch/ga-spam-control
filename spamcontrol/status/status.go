// Package status defines the available spam filter statuses.
package status

const (

	// NotSet is the default state and indicates an undefined state
	NotSet Status = 1 + iota

	// UpToDate indicates that no update is required.
	UpToDate

	// Outdated indicates that an update is required.
	Outdated

	// NotInstalled indicates that a filter is not installed.
	NotInstalled

	// Obsolete indicates that the a filter can be removed.
	Obsolete
)

var labels = []string{
	"not-set",
	"up-to-date",
	"outdated",
	"not-installed",
	"obsolete",
}

// Status defines filter statuses.
type Status int

func (status Status) String() string {
	return labels[status-1]
}
