package api

import (
	"github.com/andreaskoch/ga-spam-control/api/apicredentials"
	"github.com/andreaskoch/ga-spam-control/api/apimodel"
	"github.com/andreaskoch/ga-spam-control/api/apiservice"
)

func New(tokenStore apicredentials.TokenStorer, clientID, clientSecret string) *API {
	return &API{}
}

type API struct {
}

func (api *API) GetAccounts() ([]apimodel.Account, error) {
	apiservice.New()

	return nil, nil
}

func (api *API) GetFilters(accountID string) ([]apimodel.Filter, error) {
	apiservice.New()

	return nil, nil
}
