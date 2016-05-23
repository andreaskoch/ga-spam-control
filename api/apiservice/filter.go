package apiservice

import (
	"bytes"
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

// FilterResults is response model for Google Analytics Filter API requests.
type FilterResults struct {
	Results
	Items []Filter `json:"items"`
}

// The FilterDetail model contains the filter match type (regex, ...)
// and the filter expression value which is used to filter out
// unwanted analytics data.
type FilterDetail struct {
	Kind            string `json:"kind"`
	Field           string `json:"field"`
	MatchType       string `json:"matchType"`
	ExpressionValue string `json:"expressionValue"`
	CaseSensitive   bool   `json:"caseSensitive"`
}

// A Filter model contains the analytics filter details
// such as the filter name and type.
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

// FilterUpdate models can be used to update existing filters.
type FilterUpdate struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	Type           string             `json:"type"`
	ExcludeDetails FilterDetailUpdate `json:"excludeDetails"`
}

// The FilterDetailUpdate model can be used to update the
// expression value of an existing filter.
type FilterDetailUpdate struct {
	Field           string `json:"field"`
	ExpressionValue string `json:"expressionValue"`
}

type filterUpdateSerializer struct{}

func (filterUpdateSerializer) Serialize(writer io.Writer, filterUpdate *FilterUpdate) error {
	bytes, err := json.MarshalIndent(filterUpdate, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (filterUpdateSerializer) Deserialize(reader io.Reader) (*FilterUpdate, error) {
	decoder := json.NewDecoder(reader)
	var filterUpdate *FilterUpdate
	err := decoder.Decode(&filterUpdate)
	return filterUpdate, err
}

func serialze(object interface{}) (string, error) {
	buffer := new(bytes.Buffer)

	bytes, err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		return "", err
	}

	buffer.Write(bytes)

	return buffer.String(), nil
}
