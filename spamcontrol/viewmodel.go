package spamcontrol

import (
	"sort"

	"github.com/andreaskoch/ga-spam-control/spamcontrol/status"
)

// A StateOverview represents the spam-control status of all accounts.
type StateOverview struct {
	OverallStatus status.Status   `json:"overallStatus"`
	Accounts      []AccountStatus `json:"accounts"`
}

// An AccountStatus represents the spam-control status
// of a specific account.
type AccountStatus struct {
	AccountID   string        `json:"accountId"`
	AccountName string        `json:"accountName"`
	Status      status.Status `json:"status"`
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

	// NumberOfEntries contains of number of spam entries for current spam domain.
	NumberOfEntries int `json:"numberOfEntries"`
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
