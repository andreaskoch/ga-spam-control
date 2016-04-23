package spamcontrol

import (
	"fmt"

	"github.com/andreaskoch/ga-spam-control/api"
)

type accountProvider interface {
	// GetAccounts returns all available api.Account models.
	GetAccounts() ([]api.Account, error)

	// GetAccount returns the account with the given accountID. Returns an error if
	// the account with the given ID was not found or cannot be accessed.
	GetAccount(accountID string) (api.Account, error)
}

type remoteAccountProvider struct {
	analyticsAPI api.AnalyticsAPI
}

// GetAccounts returns all available api.Account models.
func (accountProvider remoteAccountProvider) GetAccounts() ([]api.Account, error) {
	accounts, accountsError := accountProvider.analyticsAPI.GetAccounts()
	if accountsError != nil {
		return nil, accountsError
	}

	return accounts, nil
}

// GetAccount returns the account with the given accountID. Returns an error if
// the account with the given ID was not found or cannot be accessed.
func (accountProvider remoteAccountProvider) GetAccount(accountID string) (api.Account, error) {
	accounts, accountsError := accountProvider.analyticsAPI.GetAccounts()
	if accountsError != nil {
		return api.Account{}, accountsError
	}

	for _, account := range accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return api.Account{}, fmt.Errorf("The account %q was not found.", accountID)
}
