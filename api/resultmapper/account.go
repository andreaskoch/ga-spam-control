package resultmapper

import (
	"github.com/andreaskoch/ga-spam-control/api/apimodel"
	"github.com/andreaskoch/ga-spam-control/api/apiservice"
)

func MapAccounts(sources []apiservice.Account) []apimodel.Account {

	accounts := make([]apimodel.Account, 0)
	for _, source := range sources {
		accounts = append(accounts, MapAccount(source))
	}

	return accounts
}

// MapAccount converts a apiservice.Account model into a apimodel.Account model.
func MapAccount(source apiservice.Account) apimodel.Account {
	return apimodel.Account{
		ID:   source.ID,
		Name: source.Name,
		Kind: source.Kind,
		Type: source.Type,
		Link: source.SelfLink,
	}
}
