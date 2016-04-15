package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/etix/stoppableListener"
	"golang.org/x/oauth2"
)

// GoogleAnalyticsClientID contains the client ID of the Google API credentials (see: https://console.developers.google.com/apis/credentials)
const GoogleAnalyticsClientID = "821429244906-8aki1tiaov6g2o7lr7elp41435adk9ge.apps.googleusercontent.com"

// GoogleAnalyticsClientSecret contains the client secret of the Google API credentials (see: https://console.developers.google.com/apis/credentials)
const GoogleAnalyticsClientSecret = "_WxLj0SpQ8HxqmOEyYDUTFzW"

// receiveAuthorizationCode return a redirect URL and a channel which
// receives the Google oAuth authorization code once the user has
// authorized the operation.
func receiveAuthorizationCode() (string, chan string) {

	authorizationCode := make(chan string, 1)
	listenAddress := "localhost:8080"
	route := "/authorizationCodeReceiver"
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
			}

			fmt.Fprintf(w, "No code received")
		})

		http.Serve(handler, nil)
	}()

	return fmt.Sprintf("http://%s%s", listenAddress, route), authorizationCode
}

// getAnalyticsClient returns a Google Analytics client instance.
func getAnalyticsClient() *http.Client {
	redirectURL, codeReceiver := receiveAuthorizationCode()

	conf := &oauth2.Config{
		ClientID:     GoogleAnalyticsClientID,
		ClientSecret: GoogleAnalyticsClientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/analytics.edit",
			"https://www.googleapis.com/auth/analytics.readonly",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)
	fmt.Println()

	// Use the authorization code that is pushed to the redirect URL.
	// NewTransportWithCode will do the handshake to retrieve
	// an access token and initiate a Transport that is
	// authorized and authenticated by the retrieved token.
	code := <-codeReceiver
	if len(code) == 0 {
		log.Fatal("No authorization code received.")
	}

	log.Printf("Authorization code received: %s\n", code)

	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(oauth2.NoContext, tok)
	return client
}

func main() {

	client := getAnalyticsClient()

	accountId := "578578"
	// uri := fmt.Sprintf("https://www.googleapis.com/analytics/v3/management/accounts/%s/filters", accountId)
	uri := fmt.Sprintf("https://www-googleapis-com-yb0hxtzk6st4.runscope.net/analytics/v3/management/accounts/%s/filters", accountId)
	response, err := client.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}
