package spamcontrol

import (
	"fmt"
	"sort"
)

// DomainUpdateType defines different update types of domain updates.
type DomainUpdateType int

const (

	// NotSet is the default domain update type
	// and indicates that no update type has been selected.
	NotSet DomainUpdateType = 1 + iota

	// Unchanged indicates that the domain name did
	// not change during a domain update.
	Unchanged

	// Added indicates a domain name that is new.
	Added

	// Removed indicates a domain names that has been
	// removed from the list after a domain update.
	Removed
)

func (updateType DomainUpdateType) String() string {
	return updateTypeLabels[updateType-1]
}

// IsUnchanged returns a flag indicating whether the current
// update type is Unchanged.
func (updateType DomainUpdateType) IsUnchanged() bool {
	return updateType == Unchanged
}

// IsAdded returns a flag indicating whether the current
// update type is Added.
func (updateType DomainUpdateType) IsAdded() bool {
	return updateType == Added
}

// IsRemoved returns a flag indicating whether the current
// update type is Removed.
func (updateType DomainUpdateType) IsRemoved() bool {
	return updateType == Removed
}

var updateTypeLabels = []string{
	"?",
	"o",
	"+",
	"-",
}

// An UpdateResult contains added and removed domain names.
type UpdateResult struct {
	Statistics DomainUpdateStatistics `json:"statistics"`
	Domains    []DomainUpdate         `json:"domains"`
}

// The DomainUpdateStatistics view model contains statistics
// (e.g. number of new domains) about domain update results.
type DomainUpdateStatistics struct {
	Unchanged int `json:"unchanged"`
	Added     int `json:"added"`
	Removed   int `json:"removed"`
	Total     int `json:"total"`
}

// A DomainUpdate contains the update type and the domain name
// and is the result of a domain update.
type DomainUpdate struct {
	UpdateType DomainUpdateType
	Domainname string
}

// domainUpdatesByName can be used to sort domainUpdates by name (ascending).
func domainUpdatesByName(domainUpdate1, domainUpdate2 DomainUpdate) bool {
	return domainUpdate1.Domainname < domainUpdate2.Domainname
}

// The SortDomainUpdateBy function sorts DomainUpdate objects.
type SortDomainUpdateBy func(domainUpdate1, domainUpdate2 DomainUpdate) bool

// Sort the given DomainUpdate objects.
func (by SortDomainUpdateBy) Sort(domainUpdates []DomainUpdate) {
	sorter := &domainUpdateSorter{
		domainUpdates: domainUpdates,
		by:            by,
	}

	sort.Sort(sorter)
}

type domainUpdateSorter struct {
	domainUpdates []DomainUpdate
	by            SortDomainUpdateBy
}

func (sorter *domainUpdateSorter) Len() int {
	return len(sorter.domainUpdates)
}

func (sorter *domainUpdateSorter) Swap(i, j int) {
	sorter.domainUpdates[i], sorter.domainUpdates[j] = sorter.domainUpdates[j], sorter.domainUpdates[i]
}

func (sorter *domainUpdateSorter) Less(i, j int) bool {
	return sorter.by(sorter.domainUpdates[i], sorter.domainUpdates[j])
}

// A StateOverview represents the spam-control status of all accounts.
type StateOverview struct {
	Accounts         []AccountStatus `json:"accounts"`
	KnownSpamDomains int             `json:"knownSpamDomains"`
}

// An AccountStatus represents the spam-control status
// of a specific account.
type AccountStatus struct {
	AccountID   string             `json:"accountId"`
	AccountName string             `json:"accountName"`
	Status      InstallationStatus `json:"status"`
}

// AnalysisResult represents the current spam status
// for a given analytics account.
type AnalysisResult struct {
	AccountID   string       `json:"accountId"`
	SpamDomains []SpamDomain `json:"spamDomains"`
}

// A SpamDomain model contains information about (referrer) spam domains.
type SpamDomain struct {
	// DomainName defines contains the domain name of the spam domain (e.g. "rank-checker.online")
	DomainName string `json:"domainName"`

	// SpamProbability contains the spam probability for this domain.
	SpamProbability float64 `json:"spamProbability"`
}

// spamDomainsByName can be used to sort spamDomains by name (ascending).
func spamDomainsByName(spamDomain1, spamDomain2 SpamDomain) bool {
	return spamDomain1.DomainName < spamDomain2.DomainName
}

// The SortSpamDomainsBy function sorts SpamDomain objects.
type SortSpamDomainsBy func(spamDomain1, spamDomain2 SpamDomain) bool

// Sort the given SpamDomain objects.
func (by SortSpamDomainsBy) Sort(spamDomains []SpamDomain) {
	sorter := &spamDomainSorter{
		spamDomains: spamDomains,
		by:          by,
	}

	sort.Sort(sorter)
}

type spamDomainSorter struct {
	spamDomains []SpamDomain
	by          SortSpamDomainsBy
}

func (sorter *spamDomainSorter) Len() int {
	return len(sorter.spamDomains)
}

func (sorter *spamDomainSorter) Swap(i, j int) {
	sorter.spamDomains[i], sorter.spamDomains[j] = sorter.spamDomains[j], sorter.spamDomains[i]
}

func (sorter *spamDomainSorter) Less(i, j int) bool {
	return sorter.by(sorter.spamDomains[i], sorter.spamDomains[j])
}

// InstallationStatus can be used to represent the spam-control
// status of an account. It indicates how many of all available
// filters are up-to-date.
type InstallationStatus struct {
	TotalFilters    int
	UpToDateFilters int
}

// String returns a percentage string representing the spam-control status in percent (e.g. "72%").
func (installationStatus InstallationStatus) String() string {
	if installationStatus.TotalFilters == 0 || installationStatus.UpToDateFilters == 0 {
		return "0%"
	}

	percentage := installationStatus.UpToDateFilters * 100.0 / installationStatus.TotalFilters
	return fmt.Sprintf("%d%%", percentage)
}

// A Table defines the columnes and values of a classical table.
type Table struct {
	ColumnNames []string
	Rows        [][]string
}

// MachineLearningModel contains all attributes for the machine-learning model.
type MachineLearningModel Table
