package apiservice

import (
	"encoding/json"
	"io"
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
	AccountID      string       `json:"accountId"`
	ParentLink     Link         `json:"parentLink"`
	ExcludeDetails FilterDetail `json:"excludeDetails"`
}
