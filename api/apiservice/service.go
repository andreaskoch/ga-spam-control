package apiservice

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

	if err := handleErrors(response); err != nil {
		return nil, err
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

	if err := handleErrors(response); err != nil {
		return nil, err
	}

	serializer := &filterResultsSerializer{}
	results, deserializeError := serializer.Deserialize(response.Body)
	if deserializeError != nil {
		return nil, deserializeError
	}

	return results.Items, nil
}

// CreateFilter creates a new filter for the given account ID.
func (service *GoogleAnalytics) CreateFilter(accountId string, filter Filter) error {

	buffer := new(bytes.Buffer)
	serializer := &filterSerializer{}
	serializeError := serializer.Serialize(buffer, &filter)
	if serializeError != nil {
		return serializeError
	}

	uri := fmt.Sprintf("https://%s/analytics/v3/management/accounts/%s/filters", service.apiHostname, accountId)
	response, err := service.client.Post(uri, "application/json; charset=UTF-8", buffer)
	if err != nil {
		return err
	}

	if err := handleErrors(response); err != nil {
		return err
	}

	return nil
}

func handleErrors(response *http.Response) error {
	if response.StatusCode != 200 {
		errorResponse, decodeError := decodeResponse(response.Body)
		if decodeError != nil {
			if body, err := ioutil.ReadAll(response.Body); err == nil {
				return fmt.Errorf("%s", body)
			}

			return decodeError
		}

		return fmt.Errorf("%s", errorResponse.Error.Message)
	}

	return nil
}
