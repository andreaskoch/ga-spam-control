package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// GoogleAnalyticsHostname contains the hostname of the Google Analytics API
// const GoogleAnalyticsHostname = "www.googleapis.com"
const GoogleAnalyticsHostname = "www-googleapis-com-yb0hxtzk6st4.runscope.net"

func main() {

	// token store
	homeDirPath, err := homedir.Dir()
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot determine the current users home direcotry. Error: %s", err))
	}

	tokenStoreFilePath := filepath.Join(homeDirPath, ".analytics")
	store := newTokenStore(tokenStoreFilePath)

	// credentials
	googleAnalyticsClientID := "821429244906-8aki1tiaov6g2o7lr7elp41435adk9ge.apps.googleusercontent.com"
	googleAnalyticsClientSecret := "_WxLj0SpQ8HxqmOEyYDUTFzW"

	// oAuth code receiver
	listenAddress := "localhost:8080"
	route := "/authorizationCodeReceiver"
	redirectURL := fmt.Sprintf("http://%s%s", listenAddress, route)

	// oAuth client config
	oAuthClientConfig := getAnalyticsClientConfig(googleAnalyticsClientID, googleAnalyticsClientSecret, redirectURL)
	// instantiate a Google Analytics client
	client, err := getAnalyticsClient(store, oAuthClientConfig, listenAddress, route)
	if err != nil {
		log.Fatal(err)
	}

	getAccounts(client)

	accountId := "578578"
	getFilters(client, accountId)
}
