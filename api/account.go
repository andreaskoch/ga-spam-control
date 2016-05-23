package api

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/andreaskoch/ga-spam-control/api/apiservice"
)

// toModelAccounts converts []apiservice.Account to []Account.
func toModelAccounts(sources []apiservice.Account) []Account {

	var accounts []Account
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

// An Account contains all parameters of an analytics account.
type Account struct {
	ID   string
	Kind string
	Name string
	Type string
	Link string
}

// accountsByID can be used to sort accounts by id (ascending).
func accountsByID(account1, account2 Account) bool {
	account1ID, parseAccount1IDError := strconv.ParseInt(account1.ID, 10, 64)
	if parseAccount1IDError != nil {
		panic(parseAccount1IDError)
	}

	account2ID, parseAccount2IDError := strconv.ParseInt(account2.ID, 10, 64)
	if parseAccount2IDError != nil {
		panic(parseAccount2IDError)
	}

	return fmt.Sprintf("%012d", int(account1ID)) < fmt.Sprintf("%012d", int(account2ID))
}

// SortAccountsBy sorts two Account models.
type SortAccountsBy func(account1, account2 Account) bool

// Sort a slice of Account models.
func (by SortAccountsBy) Sort(accounts []Account) {
	sorter := &accountSorter{
		accounts: accounts,
		by:       by,
	}

	sort.Sort(sorter)
}

type accountSorter struct {
	accounts []Account
	by       SortAccountsBy
}

func (sorter *accountSorter) Len() int {
	return len(sorter.accounts)
}

func (sorter *accountSorter) Swap(i, j int) {
	sorter.accounts[i], sorter.accounts[j] = sorter.accounts[j], sorter.accounts[i]
}

func (sorter *accountSorter) Less(i, j int) bool {
	return sorter.by(sorter.accounts[i], sorter.accounts[j])
}
