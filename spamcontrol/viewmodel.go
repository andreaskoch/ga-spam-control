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
	AccountID string   `json:"accountId"`
	Domains   []Domain `json:"domains"`
}

// A Domain model contains information about (referrer) spam domains.
type Domain struct {
	// DomainName defines contains the domain name of the spam domain (e.g. "rank-checker.online")
	DomainName string `json:"domainName"`

	// IsSpam contains a flag indicating whether this domain is spam or not.
	IsSpam bool `json:"isSpam"`
}

// referrerDomainsByName can be used to sort referrerDomains by name (ascending).
func referrerDomainsByName(referrerDomain1, referrerDomain2 Domain) bool {
	return referrerDomain1.DomainName < referrerDomain2.DomainName
}

// The SortDomainsBy function sorts Domain objects.
type SortDomainsBy func(referrerDomain1, referrerDomain2 Domain) bool

// Sort the given Domain objects.
func (by SortDomainsBy) Sort(referrerDomains []Domain) {
	sorter := &referrerDomainSorter{
		referrerDomains: referrerDomains,
		by:              by,
	}

	sort.Sort(sorter)
}

type referrerDomainSorter struct {
	referrerDomains []Domain
	by              SortDomainsBy
}

func (sorter *referrerDomainSorter) Len() int {
	return len(sorter.referrerDomains)
}

func (sorter *referrerDomainSorter) Swap(i, j int) {
	sorter.referrerDomains[i], sorter.referrerDomains[j] = sorter.referrerDomains[j], sorter.referrerDomains[i]
}

func (sorter *referrerDomainSorter) Less(i, j int) bool {
	return sorter.by(sorter.referrerDomains[i], sorter.referrerDomains[j])
}

// InstallationStatus can be used to represent the spam-control
// status of an account. It indicates how many of all available
// filters are up-to-date.
type InstallationStatus struct {
	TotalDomains   int
	DomainsCovered int
}

// String returns a percentage string representing the spam-control status in percent (e.g. "72%").
func (installationStatus InstallationStatus) String() string {
	if installationStatus.TotalDomains == 0 || installationStatus.DomainsCovered == 0 {
		return "0%"
	}

	percentage := installationStatus.DomainsCovered * 100.0 / installationStatus.TotalDomains
	return fmt.Sprintf("%d%%", percentage)
}
