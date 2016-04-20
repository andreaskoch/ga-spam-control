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
	response, requestError := service.client.Get(uri)
	if requestError != nil {
		return nil, fmt.Errorf("The GET request against %q failed: %s", uri, requestError.Error())
	}

	if err := handleErrors(response); err != nil {
		return nil, fmt.Errorf("The GET request against %q did not succeed: %s", uri, err.Error())
	}

	serializer := &accountResultsSerializer{}
	results, deserializeError := serializer.Deserialize(response.Body)
	if deserializeError != nil {
		return nil, fmt.Errorf("The accounts response could not be deserialized: %s", deserializeError.Error())
	}

	return results.Items, nil
}

// GetFilters returns all filters for the account with the given account ID.
func (service *GoogleAnalytics) GetFilters(accountId string) ([]Filter, error) {

	uri := fmt.Sprintf("https://%s/analytics/v3/management/accounts/%s/filters", service.apiHostname, accountId)
	response, requestError := service.client.Get(uri)
	if requestError != nil {
		return nil, fmt.Errorf("The GET request against %q failed: %s", uri, requestError.Error())
	}

	if err := handleErrors(response); err != nil {
		return nil, fmt.Errorf("The GET request against %q did not succeed: %s", uri, err.Error())
	}

	serializer := &filterResultsSerializer{}
	results, deserializeError := serializer.Deserialize(response.Body)
	if deserializeError != nil {
		return nil, fmt.Errorf("The filters response could not be deserialized: %s", deserializeError.Error())
	}

	return results.Items, nil
}

// GetProfileUserLinks returns all profile user links for the account with the given account ID.
func (service *GoogleAnalytics) GetProfileUserLinks(accountId string) ([]Profile, error) {

	uri := fmt.Sprintf("https://%s/analytics/v3/management/accounts/%s/webproperties/%s/profiles/%s/entityUserLinks", service.apiHostname, accountId, "~all", "~all")
	response, requestError := service.client.Get(uri)
	if requestError != nil {
		return nil, fmt.Errorf("The GET request against %q failed: %s", uri, requestError.Error())
	}

	if err := handleErrors(response); err != nil {
		return nil, fmt.Errorf("The GET request against %q did not succeed: %s", uri, err.Error())
	}

	serializer := &profileResultsSerializer{}
	results, deserializeError := serializer.Deserialize(response.Body)
	if deserializeError != nil {
		return nil, fmt.Errorf("The filters response could not be deserialized: %s", deserializeError.Error())
	}

	return results.Items, nil
}

// CreateFilter creates a new filter for the given account ID.
func (service *GoogleAnalytics) CreateFilter(accountId string, filter Filter) error {

	buffer := new(bytes.Buffer)
	serializer := &filterSerializer{}
	serializeError := serializer.Serialize(buffer, &filter)
	if serializeError != nil {
		return fmt.Errorf("The given filter model could not be serialized: %s", serializeError.Error())
	}

	uri := fmt.Sprintf("https://%s/analytics/v3/management/accounts/%s/filters", service.apiHostname, accountId)
	response, requestError := service.client.Post(uri, "application/json; charset=UTF-8", buffer)
	if requestError != nil {
		return fmt.Errorf("The POST request against %q failed: %s", uri, requestError.Error())
	}

	if err := handleErrors(response); err != nil {
		return fmt.Errorf("The POST request against %q did not succeed: %s", uri, err.Error())
	}

	return nil
}

// RemoveFilter deletes the given filter from the specified account.
func (service *GoogleAnalytics) RemoveFilter(accountID, filterID string) error {

	uri := fmt.Sprintf("https://%s/analytics/v3/management/accounts/%s/filters/%s", service.apiHostname, accountID, filterID)

	request, createRequestError := http.NewRequest(http.MethodDelete, uri, nil)
	if createRequestError != nil {
		return fmt.Errorf("The DELETE request could not be created: %s", createRequestError.Error())
	}

	response, requestError := service.client.Do(request)
	if requestError != nil {
		return fmt.Errorf("The DELETE request against %q failed: %s", uri, requestError.Error())
	}

	if err := handleErrors(response); err != nil {
		return fmt.Errorf("The DELETE request against %q did not succeed: %s", uri, err.Error())
	}

	return nil
}

func handleErrors(response *http.Response) error {
	if response.StatusCode != 200 && response.StatusCode != 204 {
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
