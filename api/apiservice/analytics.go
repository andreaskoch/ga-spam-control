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

type AnalyticsDataResults struct {
	Results
	Data AnalyticsData `json:"dataTable"`
}

type AnalyticsData struct {
	Cols []TableColumn `json:"cols"`
	Rows []TableRow    `json:"rows"`
}

type TableColumn struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type TableRow struct {
	Cell []TableCell `json:"c"`
}

type TableCell struct {
	Value string `json:"v"`
}
