package apiservice

import (
	"fmt"
	"net/http"
)

// GoogleAnalyticsHostname contains the hostname of the Google Analytics API
// const GoogleAnalyticsHostname = "www.googleapis.com"
const GoogleAnalyticsHostname = "www-googleapis-com-yb0hxtzk6st4.runscope.net"

func New() *GoogleAnalytics {
	return &GoogleAnalytics{
		apiHostname: GoogleAnalyticsHostname,
	}
}

type GoogleAnalytics struct {
	apiHostname string
	client      *http.Client
}

// GetAccounts returns all accessible accounts from the given API client.
func (service *GoogleAnalytics) GetAccounts() ([]Account, error) {

	uri := fmt.Sprintf("https://%s/analytics/v3/management/accounts", service.apiHostname)
	response, apiError := service.client.Get(uri)
	if apiError != nil {
		return nil, apiError
	}

	serializer := &accountResultsSerializer{}
	results, deserializeError := serializer.Deserialize(response.Body)
	if deserializeError != nil {
		return nil, deserializeError
	}

	return results.Items, nil
}

// GetFilters returns all filters for the account with the given account ID.
func (service *GoogleAnalytics) GetFilters(accountId string) ([]Filter, error) {

	uri := fmt.Sprintf("https://%s/analytics/v3/management/accounts/%s/filters", service.apiHostname, accountId)
	response, err := service.client.Get(uri)
	if err != nil {
		return nil, err
	}

	fmt.Println(response)
	return nil, nil
}
