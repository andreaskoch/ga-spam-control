package spamcontrol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/andreaskoch/ga-spam-control/api"
)

// The SpamDetector interface provides a functions detecting spam
// in analytics data.
type SpamDetector interface {
	// GetSpamRating returns the rated spam score for the given analytics data.
	GetSpamRating(analyticsData []api.AnalyticsDataRow) (RatedAnalyticsData, error)
}

// NewDetector create a new SpamDetector instance.
func NewDetector() SpamDetector {
	return azureMLSpamDetection{}
}

// azureMLSpamDetection uses Azure ML Studio to detect spam in analytics data.
type azureMLSpamDetection struct {
}

// GetSpamRating returns the rated spam score for the given analytics data.
func (spamDetection azureMLSpamDetection) GetSpamRating(analyticsData []api.AnalyticsDataRow) (RatedAnalyticsData, error) {

	inputSerializer := &inputRequestSerializer{}
	outputSerializer := &spamScoreResponseSerializer{}

	transformedData := analyticsDataToMachineLearningModel(analyticsData)
	inputRequest := rowsToInputRequest(transformedData)

	buffer := new(bytes.Buffer)
	serializeError := inputSerializer.Serialize(buffer, &inputRequest)
	if serializeError != nil {
		return nil, fmt.Errorf("The input request model could not be serialized: %s", serializeError.Error())
	}

	// uri := "https://europewest.services.azureml.net/workspaces/7cd19bff3eb34765a374b73e9820efba/services/7b8554f367af497382f5fde320121321/execute?api-version=2.0&details=true"
	uri := "https://europewest-services-azureml-net-yb0hxtzk6st4.runscope.net/workspaces/7cd19bff3eb34765a374b73e9820efba/services/7b8554f367af497382f5fde320121321/execute?api-version=2.0&details=true"
	request, createRequestError := http.NewRequest(http.MethodPost, uri, buffer)
	if createRequestError != nil {
		return nil, fmt.Errorf("The PUT request could not be created: %s", createRequestError.Error())
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "x4DKP9sdRUpq1NkWQ0EgIBgtXjbBtxOxqgFYTSsznZO/XTqT4XHJNGMuwtzmEGHNThkh2CjYQ6gb/AkoHYzmuw=="))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, requestError := http.DefaultClient.Do(request)
	if requestError != nil {
		return nil, requestError
	}

	spamScoreResponse, deserializeError := outputSerializer.Deserialize(response.Body)
	if deserializeError != nil {
		return nil, deserializeError
	}

	ratedAnalyticsData, err := spamScoreResponseToRatedAnalyticsData(analyticsData, *spamScoreResponse)
	if err != nil {
		return nil, err
	}

	return ratedAnalyticsData, nil
}

type spamScoreResponseSerializer struct{}

func (spamScoreResponseSerializer) Serialize(writer io.Writer, spamScoreResponse *spamScoreResponse) error {
	bytes, err := json.MarshalIndent(spamScoreResponse, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (spamScoreResponseSerializer) Deserialize(reader io.Reader) (*spamScoreResponse, error) {
	decoder := json.NewDecoder(reader)
	var spamScoreResponse *spamScoreResponse
	err := decoder.Decode(&spamScoreResponse)
	return spamScoreResponse, err
}

type spamScoreValue struct {
	ColumnNames []string   `json:"ColumnNames"`
	ColumnTypes []string   `json:"ColumnTypes"`
	Values      [][]string `json:"Values"`
}

type spamScore struct {
	Type  string         `json:"type"`
	Value spamScoreValue `json:"value"`
}

type results struct {
	SpamScore spamScore `json:"spamScore"`
}

type spamScoreResponse struct {
	Results results `json:"Results"`
}

// RatedAnalyticsDataRow contains a rated analyitcs data row.
type RatedAnalyticsDataRow struct {
	api.AnalyticsDataRow

	IsSpam      bool
	Probability float64
}

// RatedAnalyticsData contains a set of rated analytics data rows
type RatedAnalyticsData []RatedAnalyticsDataRow

func spamScoreResponseToRatedAnalyticsData(requestValues []api.AnalyticsDataRow, response spamScoreResponse) (RatedAnalyticsData, error) {

	responseValues := response.Results.SpamScore.Value.Values
	// if len(responseValues) != len(requestValues) {
	// return RatedAnalyticsData{}, fmt.Errorf("Response size does not match request size.")
	// }

	results := make(RatedAnalyticsData, 0)

	for index, spamScore := range responseValues {

		// look up the source/domain from the request data
		requestData := requestValues[index]
		dataRow := api.AnalyticsDataRow{
			Source: requestData.Source,
		}

		isSpam, err := strconv.ParseBool(spamScore[0])
		if err != nil {
			return nil, fmt.Errorf("Unable to parse %q: %s", spamScore[0], err.Error())
		}

		propability, err := strconv.ParseFloat(spamScore[1], 64)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse %q: %s", spamScore[1], err.Error())
		}

		row := RatedAnalyticsDataRow{dataRow, isSpam, propability}
		results = append(results, row)

	}

	return results, nil
}

func rowsToInputRequest(analyticsData Table) inputRequest {

	request := inputRequest{
		Inputs: inputs{
			ReferrerData: referrerData{
				ColumnNames: analyticsData.ColumnNames,
				Values:      analyticsData.Rows,
			},
		},
	}

	return request
}

type referrerData struct {
	ColumnNames []string   `json:"ColumnNames"`
	Values      [][]string `json:"Values"`
}

type inputs struct {
	ReferrerData referrerData `json:"referrerData"`
}

type inputRequestParameters struct {
}

type inputRequest struct {
	Inputs           inputs                 `json:"Inputs"`
	GlobalParameters inputRequestParameters `json:"GlobalParameters"`
}

type inputRequestSerializer struct{}

func (inputRequestSerializer) Serialize(writer io.Writer, inputRequest *inputRequest) error {
	bytes, err := json.MarshalIndent(inputRequest, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (inputRequestSerializer) Deserialize(reader io.Reader) (*inputRequest, error) {
	decoder := json.NewDecoder(reader)
	var inputRequest *inputRequest
	err := decoder.Decode(&inputRequest)
	return inputRequest, err
}
