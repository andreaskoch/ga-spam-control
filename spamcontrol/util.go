package spamcontrol

import (
	"strconv"

	"github.com/andreaskoch/ga-spam-control/api"
)

// unique returns a copy of the given string array
// cleaned from duplicate entries.
func unique(list []string) []string {
	var cleanedList []string

	index := make(map[string]int)
	for _, entry := range list {
		if _, exists := index[entry]; exists {
			continue
		}

		index[entry] = 1

		cleanedList = append(cleanedList, entry)
	}

	return cleanedList
}

// A Table defines the columnes and values of a classical table.
type Table struct {
	ColumnNames []string
	Rows        [][]string
}

// TrainingData contains all attributes for the machine-learning model.
type TrainingData Table

const trainingdataFalse = "0"
const trainingdataTrue = "1"
const trainingdataNotset = "(not set)"
const trainingdataNewvisitor = "(New Visitor)"
const trainingdataDirect = "(direct)"

func analyticsDataToTrainingData(rows []api.AnalyticsDataRow) Table {

	var values [][]string
	for _, row := range rows {

		isNewVisitor := trainingdataFalse
		if row.UserType == trainingdataNewvisitor {
			isNewVisitor = trainingdataTrue
		}

		fullReferrerIsSet := trainingdataTrue
		if row.FullReferrer == trainingdataDirect {
			fullReferrerIsSet = trainingdataFalse
		}

		mediumIsSet := trainingdataTrue
		if row.Medium == trainingdataNotset {
			mediumIsSet = trainingdataFalse
		}

		networkDomainIsSet := trainingdataTrue
		if row.NetworkDomain == trainingdataNotset {
			networkDomainIsSet = trainingdataFalse
		}

		networkLocationIsSet := trainingdataTrue
		if row.NetworkLocation == trainingdataNotset {
			networkLocationIsSet = trainingdataFalse
		}

		landingPagePathIsSet := trainingdataTrue
		if row.LandingPagePath == "/" {
			landingPagePathIsSet = trainingdataFalse
		}

		rowValues := []string{
			isNewVisitor,
			fullReferrerIsSet,
			row.Source,
			mediumIsSet,
			networkDomainIsSet,
			networkLocationIsSet,
			landingPagePathIsSet,
			strconv.FormatInt(row.Sessions, 10),
			strconv.FormatFloat(row.BounceRate, 'f', -1, 32),
			strconv.FormatFloat(row.PageviewsPerSession, 'f', -1, 32),
			strconv.FormatFloat(row.TimeOnPage, 'f', -1, 32),
		}

		values = append(values, rowValues)
	}

	request := Table{
		ColumnNames: []string{
			"isNewVisitor",
			"fullReferrerIsSet",
			"ga:source",
			"mediumIsSet",
			"networkDomainIsSet",
			"networkLocationIsSet",
			"landingPagePathIsSet",
			"ga:sessions",
			"ga:bounceRate",
			"ga:pageviewsPerSession",
			"ga:timeOnPage",
		},
		Rows: values,
	}

	return request
}
