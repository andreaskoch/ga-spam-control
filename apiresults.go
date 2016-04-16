package main

import (
	"encoding/json"
	"io"
	"time"
)

type accountResultsSerializer struct{}

func (accountResultsSerializer) Serialize(writer io.Writer, accountResults *AccountResults) error {
	bytes, err := json.MarshalIndent(accountResults, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (accountResultsSerializer) Deserialize(reader io.Reader) (*AccountResults, error) {
	decoder := json.NewDecoder(reader)
	var accountResults *AccountResults
	err := decoder.Decode(&accountResults)
	return accountResults, err
}

type Results struct {
	Kind         string `json:"kind"`
	Username     string `json:"username"`
	TotalResults int    `json:"totalResults"`
	StartIndex   int    `json:"startIndex"`
	ItemsPerPage int    `json:"itemsPerPage"`
}

type AccountResults struct {
	Results
	Items []Account `json:"items"`
}

type AccountPermissions struct {
	Effective []string `json:"effective"`
}

type Account struct {
	ID          string             `json:"id"`
	Kind        string             `json:"kind"`
	SelfLink    string             `json:"selfLink"`
	Name        string             `json:"name"`
	Permissions AccountPermissions `json:"permissions"`
	Created     time.Time          `json:"created"`
	Updated     time.Time          `json:"updated"`
	ChildLink   Link               `json:"childLink"`
}

type FilterResults struct {
	Results
	Items []Filter `json:"items"`
}

type Link struct {
	Type string `json:"type"`
	Href string `json:"href"`
}

type FilterDetail struct {
	Kind            string `json:"kind"`
	Field           string `json:"field"`
	MatchType       string `json:"matchType"`
	ExpressionValue string `json:"expressionValue"`
	CaseSensitive   bool   `json:"caseSensitive"`
}

type Filter struct {
	ID             string       `json:"id"`
	Kind           string       `json:"kind"`
	SelfLink       string       `json:"selfLink"`
	AccountID      string       `json:"accountId"`
	Name           string       `json:"name"`
	Type           string       `json:"type"`
	Created        time.Time    `json:"created"`
	Updated        time.Time    `json:"updated"`
	ParentLink     Link         `json:"parentLink"`
	ExcludeDetails FilterDetail `json:"excludeDetails"`
}
