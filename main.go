package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

func receiveAuthorizationCode() (string, chan string) {

	authorizationCode := make(chan string, 1)
	listenAddress := "localhost:8080"
	route := "/authorizationCodeReceiver"
	go func() {
		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if len(code) > 0 {
				fmt.Fprintf(w, "%s", code)
				authorizationCode <- code
				return
				// panic(fmt.Sprintf("Code received: %s", code))
			}

			fmt.Fprintf(w, "No code received")
		})

		log.Fatal(http.ListenAndServe(listenAddress, nil))
	}()

	return fmt.Sprintf("http://%s%s", listenAddress, route), authorizationCode
}

func main() {

	accountId := "578578"
	redirectURL, codeReceiver := receiveAuthorizationCode()

	conf := &oauth2.Config{
		ClientID:     "821429244906-8aki1tiaov6g2o7lr7elp41435adk9ge.apps.googleusercontent.com",
		ClientSecret: "_WxLj0SpQ8HxqmOEyYDUTFzW",
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
		fmt.Fatal("No authorization code received.")
	}

	log.Printf("Authorization code received: %s\n", code)

	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(oauth2.NoContext, tok)
	// uri := fmt.Sprintf("https://www.googleapis.com/analytics/v3/management/accounts/%s/filters", accountId)
	uri := fmt.Sprintf("https://www-googleapis-com-yb0hxtzk6st4.runscope.net/analytics/v3/management/accounts/%s/filters", accountId)
	response, err := client.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}
