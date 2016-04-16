package cli

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/andreaskoch/ga-spam-control/api"
	"github.com/andreaskoch/ga-spam-control/cli/credentials"
	"github.com/mitchellh/go-homedir"
)

func main() {

	// create a token store
	homeDirPath, err := homedir.Dir()
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot determine the current users home direcotry. Error: %s", err))
	}

	tokenStoreFilePath := filepath.Join(homeDirPath, ".analytics")
	tokenStore := credentials.NewTokenStore(tokenStoreFilePath)

	// create a new analytis API instance
	googleAnalyticsClientID := "821429244906-8aki1tiaov6g2o7lr7elp41435adk9ge.apps.googleusercontent.com"
	googleAnalyticsClientSecret := "_WxLj0SpQ8HxqmOEyYDUTFzW"
	analyticsAPI := api.New(tokenStore, googleAnalyticsClientID, googleAnalyticsClientSecret)

	// get all available accounts
	accounts, accountsError := analyticsAPI.GetAccounts()
	if accountsError != nil {
		log.Fatal(accountsError)
	}

	for _, account := range accounts {

		// get all filters for account
		filters, filtersError := analyticsAPI.GetFilters(account.ID)
		if filtersError != nil {
			log.Fatal(filtersError)
		}

		for _, filter := range filters {
			log.Printf("%#v\n", filter)
		}

	}
}
