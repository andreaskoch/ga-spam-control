package spamcontrol

import "github.com/andreaskoch/ga-spam-control/api"

type accountProvider interface {
	// GetAccounts returns all available api.Account models.
	GetAccounts() ([]api.Account, error)
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
