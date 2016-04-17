package api

import (
	"github.com/andreaskoch/ga-spam-control/api/apiservice"
)

// toModelAccounts converts []apiservice.Account to []Account.
func toModelAccounts(sources []apiservice.Account) []Account {

	accounts := make([]Account, 0)
	for _, source := range sources {
		accounts = append(accounts, toModelAccount(source))
	}

	return accounts
}

// toModelAccount converts a apiservice.Account into a Account.
func toModelAccount(source apiservice.Account) Account {
	return Account{
		ID:   source.ID,
		Name: source.Name,
		Kind: source.Kind,
		Type: source.Type,
		Link: source.SelfLink,
	}
}

type Account struct {
	ID   string
	Kind string
	Name string
	Type string
	Link string
}
