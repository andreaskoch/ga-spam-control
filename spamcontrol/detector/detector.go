// Package detector contains a spam detector implementation that uses
// Azure Machine Learning web services to detect spam in analytics data.
package detector

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
	GetSpamRating(analyticsData api.AnalyticsData) (RatedAnalyticsData, error)
}

// New create a new SpamDetector instance.
func New() SpamDetector {
	return azureMLSpamDetection{}
}

// azureMLSpamDetection uses Azure ML Studio to detect spam in analytics data.
type azureMLSpamDetection struct {
}

// GetSpamRating returns the rated spam score for the given analytics data.
func (spamDetection azureMLSpamDetection) GetSpamRating(analyticsData api.AnalyticsData) (RatedAnalyticsData, error) {

	inputSerializer := &inputRequestSerializer{}
	outputSerializer := &spamScoreResponseSerializer{}

	inputRequest := rowsToInputRequest(analyticsData)

	buffer := new(bytes.Buffer)
	serializeError := inputSerializer.Serialize(buffer, &inputRequest)
	if serializeError != nil {
		return nil, fmt.Errorf("The input request model could not be serialized: %s", serializeError.Error())
	}

	uri := "https://europewest.services.azureml.net/workspaces/7cd19bff3eb34765a374b73e9820efba/services/9def2d38e1d6422c8115fcb04bec56c7/execute?api-version=2.0&details=true"
	// uri := "https://europewest-services-azureml-net-yb0hxtzk6st4.runscope.net/workspaces/7cd19bff3eb34765a374b73e9820efba/services/9def2d38e1d6422c8115fcb04bec56c7/execute?api-version=2.0&details=true"
	request, createRequestError := http.NewRequest(http.MethodPost, uri, buffer)
	if createRequestError != nil {
		return nil, fmt.Errorf("The PUT request could not be created: %s", createRequestError.Error())
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "1UC1zBQmL63esM1NH6+udfmnmCJTvysNrIJ4HC1DwHFpSkJ2oKfO6bhRBErRjyuRaS6Kkq/FaGR4MENcdEtQgQ=="))
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

	ratedAnalyticsData, err := spamScoreResponseToRatedAnalyticsData(*spamScoreResponse)
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

func spamScoreResponseToRatedAnalyticsData(response spamScoreResponse) (RatedAnalyticsData, error) {

	results := make(RatedAnalyticsData, 0)

	for _, spamScore := range response.Results.SpamScore.Value.Values {

		sessions, err := strconv.ParseInt(spamScore[7], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse %q: %s", spamScore[7], err.Error())
		}

		bounceRate, err := strconv.ParseFloat(spamScore[8], 64)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse %q: %s", spamScore[8], err.Error())
		}

		pageViewsPerSession, err := strconv.ParseFloat(spamScore[9], 64)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse %q: %s", spamScore[9], err.Error())
		}

		timeOnPage, err := strconv.ParseFloat(spamScore[10], 64)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse %q: %s", spamScore[10], err.Error())
		}

		dataRow := api.AnalyticsDataRow{
			UserType:            spamScore[0],
			FullReferrer:        spamScore[1],
			Source:              spamScore[2],
			Medium:              spamScore[3],
			NetworkDomain:       spamScore[4],
			NetworkLocation:     spamScore[5],
			LandingPagePath:     spamScore[6],
			Sessions:            sessions,
			BounceRate:          bounceRate,
			PageviewsPerSession: pageViewsPerSession,
			TimeOnPage:          timeOnPage,
		}

		isSpam, err := strconv.ParseBool(spamScore[11])
		if err != nil {
			return nil, fmt.Errorf("Unable to parse %q: %s", spamScore[11], err.Error())
		}

		propability, err := strconv.ParseFloat(spamScore[12], 64)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse %q: %s", spamScore[12], err.Error())
		}

		row := RatedAnalyticsDataRow{dataRow, isSpam, propability}
		results = append(results, row)

	}

	return results, nil
}

func rowsToInputRequest(rows []api.AnalyticsDataRow) inputRequest {

	var values [][]string
	for _, row := range rows {
		rowValues := []string{
			row.UserType,
			row.FullReferrer,
			row.Source,
			row.Medium,
			row.NetworkDomain,
			row.NetworkLocation,
			row.LandingPagePath,
			strconv.FormatInt(row.Sessions, 10),
			strconv.FormatFloat(row.BounceRate, 'f', -1, 32),
			strconv.FormatFloat(row.PageviewsPerSession, 'f', -1, 32),
			strconv.FormatFloat(row.TimeOnPage, 'f', -1, 32),
		}

		values = append(values, rowValues)
	}

	request := inputRequest{
		Inputs: inputs{
			ReferrerData: referrerData{
				ColumnNames: []string{
					"ga:userType",
					"ga:fullReferrer",
					"ga:source",
					"ga:medium",
					"ga:networkDomain",
					"ga:networkLocation",
					"ga:landingPagePath",
					"ga:sessions",
					"ga:bounceRate",
					"ga:pageviewsPerSession",
					"ga:timeOnPage",
				},
				Values: values,
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
