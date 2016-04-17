package resultmapper

import (
	"github.com/andreaskoch/ga-spam-control/api/apimodel"
	"github.com/andreaskoch/ga-spam-control/api/apiservice"
)

func ToModelAccounts(sources []apiservice.Account) []apimodel.Account {

	accounts := make([]apimodel.Account, 0)
	for _, source := range sources {
		accounts = append(accounts, ToModelAccount(source))
	}

	return accounts
}

// ToModelAccount converts a apiservice.Account model into a apimodel.Account model.
func ToModelAccount(source apiservice.Account) apimodel.Account {
	return apimodel.Account{
		ID:   source.ID,
		Name: source.Name,
		Kind: source.Kind,
		Type: source.Type,
		Link: source.SelfLink,
	}
}
