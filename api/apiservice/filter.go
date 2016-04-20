package apiservice

import (
	"encoding/json"
	"io"
	"time"
)

type filterResultsSerializer struct{}

func (filterResultsSerializer) Serialize(writer io.Writer, filterResults *FilterResults) error {
	bytes, err := json.MarshalIndent(filterResults, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (filterResultsSerializer) Deserialize(reader io.Reader) (*FilterResults, error) {
	decoder := json.NewDecoder(reader)
	var filterResults *FilterResults
	err := decoder.Decode(&filterResults)
	return filterResults, err
}

type filterSerializer struct{}

func (filterSerializer) Serialize(writer io.Writer, filter *Filter) error {
	bytes, err := json.MarshalIndent(filter, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (filterSerializer) Deserialize(reader io.Reader) (*Filter, error) {
	decoder := json.NewDecoder(reader)
	var filter *Filter
	err := decoder.Decode(&filter)
	return filter, err
}

type FilterResults struct {
	Results
	Items []Filter `json:"items"`
}

type FilterDetail struct {
	Kind            string `json:"kind"`
	Field           string `json:"field"`
	MatchType       string `json:"matchType"`
	ExpressionValue string `json:"expressionValue"`
	CaseSensitive   bool   `json:"caseSensitive"`
}

type Filter struct {
	Item
	Name           string       `json:"name"`
	Type           string       `json:"type"`
	Created        time.Time    `json:"created"`
	Updated        time.Time    `json:"updated"`
	AccountID      string       `json:"accountId"`
	ParentLink     Link         `json:"parentLink"`
	ExcludeDetails FilterDetail `json:"excludeDetails"`
}
