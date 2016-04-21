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
