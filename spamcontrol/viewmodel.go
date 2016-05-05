package spamcontrol

import "github.com/andreaskoch/ga-spam-control/spamcontrol/status"

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
