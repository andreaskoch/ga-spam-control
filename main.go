package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/etix/stoppableListener"
	"golang.org/x/oauth2"
)

// GoogleAnalyticsHostname contains the hostname of the Google Analytics API
// const GoogleAnalyticsHostname = "www.googleapis.com"
const GoogleAnalyticsHostname = "www-googleapis-com-yb0hxtzk6st4.runscope.net"

// getAnalyticsClientConfig returns the oAuth client configuration for the Google Analytics API.
func getAnalyticsClientConfig(clientId, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/analytics",
			"https://www.googleapis.com/auth/analytics.edit",
			"https://www.googleapis.com/auth/analytics.readonly",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
}

// receiveAuthorizationCode returns the authorization code.
func receiveAuthorizationCode(conf *oauth2.Config, listenAddress, route string) (string, error) {

	authorizationCode := make(chan string, 1)

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)
	fmt.Println()
	go func() {

		listener, err := net.Listen("tcp", listenAddress)
		if err != nil {
			log.Fatal(err)
		}

		handler := stoppableListener.Handle(listener)

		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if len(code) > 0 {
				fmt.Fprintf(w, "%s", code)
				authorizationCode <- code

				handler.Stop <- true
				return
			}

			fmt.Fprintf(w, "No code received")
		})

		http.Serve(handler, nil)
	}()

	// Use the authorization code that is pushed to the redirect URL.
	// NewTransportWithCode will do the handshake to retrieve
	// an access token and initiate a Transport that is
	// authorized and authenticated by the retrieved token.
	code := <-authorizationCode
	if len(code) == 0 {
		return "", fmt.Errorf("No authorization code received.")
	}

	return code, nil
}

// getAnalyticsClient returns a Google Analytics client instance.
func getAnalyticsClient(oAuthClientConfig *oauth2.Config, listenAddress, route string) *http.Client {
	code, err := receiveAuthorizationCode(oAuthClientConfig, listenAddress, route)
	if err != nil {
		log.Fatal(err)
	}

	tok, err := oAuthClientConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}

	client := oAuthClientConfig.Client(oauth2.NoContext, tok)
	return client
}

func main() {

	// credentials
	googleAnalyticsClientID := "821429244906-8aki1tiaov6g2o7lr7elp41435adk9ge.apps.googleusercontent.com"
	googleAnalyticsClientSecret := "_WxLj0SpQ8HxqmOEyYDUTFzW"

	// oAuth code receiver
	listenAddress := "localhost:8080"
	route := "/authorizationCodeReceiver"
	redirectURL := fmt.Sprintf("http://%s%s", listenAddress, route)

	// oAuth client config
	oAuthClientConfig := getAnalyticsClientConfig(googleAnalyticsClientID, googleAnalyticsClientSecret, redirectURL)
	client := getAnalyticsClient(oAuthClientConfig, listenAddress, route)
	getAccounts(client)

	accountId := "578578"
	getFilters(client, accountId)
}

// getAccounts returns all accessible accounts.
func getAccounts(apiClient *http.Client) error {

	uri := fmt.Sprintf("https://%s/analytics/v3/management/accounts", GoogleAnalyticsHostname)
	response, err := apiClient.Get(uri)
	if err != nil {
		return err
	}

	fmt.Println(response)
	return nil
}

// getFilters returns all filters for the account with the given account ID.
func getFilters(apiClient *http.Client, accountId string) error {

	uri := fmt.Sprintf("https://%s/analytics/v3/management/accounts/%s/filters", GoogleAnalyticsHostname, accountId)
	response, err := apiClient.Get(uri)
	if err != nil {
		return err
	}

	fmt.Println(response)
	return nil
}
