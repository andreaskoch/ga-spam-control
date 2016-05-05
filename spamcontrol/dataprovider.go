package spamcontrol

import "github.com/andreaskoch/ga-spam-control/api"

type analyticsDataProvider interface {
	// GetAnalyticsData returns the api.AnalyticsData for the given account.
	GetAnalyticsData(accountID string) (api.AnalyticsData, error)
}

type remoteAnalyticsDataProvider struct {
	analyticsAPI api.AnalyticsAPI
}

// GetAnalyticsData returns the api.AnalyticsData for the given account.
func (analyticsProvider *remoteAnalyticsDataProvider) GetAnalyticsData(accountID string) (api.AnalyticsData, error) {

	analyticsData, analyticsDataError := analyticsProvider.analyticsAPI.GetAnalyticsData(accountID)
	if analyticsDataError != nil {
		return api.AnalyticsData{}, analyticsDataError
	}

	return analyticsData, nil

}
