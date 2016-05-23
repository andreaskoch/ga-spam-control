package apiservice

import (
	"encoding/json"
	"io"
)

type analyticsDataResultsSerializer struct{}

func (analyticsDataResultsSerializer) Serialize(writer io.Writer, analyticsDataResults *AnalyticsDataResults) error {
	bytes, err := json.MarshalIndent(analyticsDataResults, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (analyticsDataResultsSerializer) Deserialize(reader io.Reader) (*AnalyticsDataResults, error) {
	decoder := json.NewDecoder(reader)
	var analyticsDataResults *AnalyticsDataResults
	err := decoder.Decode(&analyticsDataResults)
	return analyticsDataResults, err
}

type analyticsDataSerializer struct{}

func (analyticsDataSerializer) Serialize(writer io.Writer, analyticsData *AnalyticsData) error {
	bytes, err := json.MarshalIndent(analyticsData, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (analyticsDataSerializer) Deserialize(reader io.Reader) (*AnalyticsData, error) {
	decoder := json.NewDecoder(reader)
	var analyticsData *AnalyticsData
	err := decoder.Decode(&analyticsData)
	return analyticsData, err
}

// AnalyticsDataResults is response model for Google Analytics data API requests.
type AnalyticsDataResults struct {
	Results
	Data AnalyticsData `json:"dataTable"`
}

// AnalyticsData represents analytics reports data in columns and rows.
type AnalyticsData struct {
	Cols []TableColumn `json:"cols"`
	Rows []TableRow    `json:"rows"`
}

// TableColumn defines analytics data table columns.
type TableColumn struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

// TableRow defines analytics data table rows.
type TableRow struct {
	Cell []TableCell `json:"c"`
}

// TableCell defines analytics data table cell/value.
type TableCell struct {
	Value string `json:"v"`
}
