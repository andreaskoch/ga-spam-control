package apiservice

import (
	"fmt"
	"net/http"

	"github.com/andreaskoch/ga-spam-control/api/apicredentials"
)

// GoogleAnalyticsHostname contains the hostname of the Google Analytics API
// const GoogleAnalyticsHostname = "www.googleapis.com"
const GoogleAnalyticsHostname = "www-googleapis-com-yb0hxtzk6st4.runscope.net"

func New(tokenStore apicredentials.TokenStorer, clientID, clientSecret string) (*GoogleAnalytics, error) {

	// oAuth code receiver
	listenAddress := "localhost:8080"
	route := "/authorizationCodeReceiver"
	redirectURL := fmt.Sprintf("http://%s%s", listenAddress, route)

	// oAuth client config
	oAuthClientConfig := getAnalyticsClientConfig(clientID, clientSecret, redirectURL)

	// instantiate a Google Analytics client
	client, err := getAnalyticsClient(tokenStore, oAuthClientConfig, listenAddress, route)
	if err != nil {
		return nil, err
	}

	return &GoogleAnalytics{
		apiHostname: GoogleAnalyticsHostname,
		client:      client,
	}, nil
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

	serializer := &filterResultsSerializer{}
	results, deserializeError := serializer.Deserialize(response.Body)
	if deserializeError != nil {
		return nil, deserializeError
	}

	return results.Items, nil
}
