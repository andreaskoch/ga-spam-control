package apiservice

import (
	"encoding/json"
	"io"
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

type AccountResults struct {
	Results
	Items []Account `json:"items"`
}

type AccountPermissions struct {
	Effective []string `json:"effective"`
}

type Account struct {
	Item
	Permissions AccountPermissions `json:"permissions"`
	ChildLink   Link               `json:"childLink"`
}
