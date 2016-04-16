package apiservice

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
