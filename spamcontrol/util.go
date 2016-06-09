package spamcontrol

import (
	"strconv"
	"strings"

	"github.com/andreaskoch/ga-spam-control/api"
)

// removeDuplicatesFromList returns a copy of the given string array
// cleaned from duplicate entries.
func removeDuplicatesFromList(list []string) []string {
	var cleanedList []string

	index := make(map[string]int)
	for _, entry := range list {
		if _, exists := index[entry]; exists {
			index[entry] = index[entry] + 1
			continue
		}

		index[entry] = 1

		cleanedList = append(cleanedList, entry)
	}

	return cleanedList
}

// removeDuplicatesFromTable returns a copy of the given table
// with all duplicate rows removed.
func removeDuplicatesFromTable(table [][]string) [][]string {
	var cleanedTable [][]string

	index := make(map[string]int)

	for _, row := range table {
		key := strings.Join(row, ",")

		if _, exists := index[key]; exists {
			index[key] = index[key] + 1
			continue
		}

		index[key] = 1

		cleanedTable = append(cleanedTable, row)
	}

	return cleanedTable
}

const trainingdataFalse = "0"
const trainingdataTrue = "1"
const trainingdataNotset = "(not set)"
const trainingdataNewvisitor = "(New Visitor)"
const trainingdataDirect = "(direct)"

func analyticsDataToMachineLearningModel(rows []api.AnalyticsDataRow) Table {

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
			mediumIsSet,
			networkDomainIsSet,
			networkLocationIsSet,
			landingPagePathIsSet,
			strconv.FormatInt(row.Sessions, 10),
			strconv.FormatFloat(row.BounceRate, 'f', -1, 32),
			strconv.FormatFloat(row.PageviewsPerSession, 'f', -1, 32),
			strconv.FormatFloat(row.TimeOnPage, 'f', -1, 32),
			row.Source,
		}

		values = append(values, rowValues)
	}

	request := Table{
		ColumnNames: []string{
			"isNewVisitor",
			"fullReferrerIsSet",
			"mediumIsSet",
			"networkDomainIsSet",
			"networkLocationIsSet",
			"landingPagePathIsSet",
			"sessions",
			"bounceRate",
			"pageviewsPerSession",
			"timeOnPage",
			"source",
		},
		Rows: values,
	}

	return request
}
