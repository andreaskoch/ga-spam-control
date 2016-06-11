package spamcontrol

import (
	"fmt"
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

func analyticsDataToMachineLearningModel(rows []api.AnalyticsDataRow) Table {

	// normalize the analytics data
	values := normalizeAnalyticsData(rows)

	// determine usage
	usageIndex := make(map[string]int)
	for _, row := range values {
		key := strings.Join(row, ",")

		if _, exists := usageIndex[key]; exists {
			usageIndex[key] = usageIndex[key] + 1
			continue
		}

		usageIndex[key] = 1
	}

	// group
	alreadySeen := make(map[string]bool)
	var groupedValues [][]string
	for _, row := range values {
		key := strings.Join(row, ",")

		// handle every key only once
		if _, exists := alreadySeen[key]; exists {
			continue
		}

		// append the number of duplicates to the row
		numberOfDuplicates := usageIndex[key]
		row = append(row, fmt.Sprintf("%v", numberOfDuplicates))

		groupedValues = append(groupedValues, row)

		// make sure we don't add this value again
		alreadySeen[key] = true
	}

	return Table{
		ColumnNames: []string{
			"isNewVisitor",
			"fullReferrerIsSet",
			"isReferral",
			"landingPagePathIsSet",
			"sessions",
			"bounceRate",
			"pageviewsPerSession",
			"timeOnPage",
			"source",
			"numberOfDuplicates",
		},
		Rows: groupedValues,
	}
}

const trainingdataFalse = "0"
const trainingdataTrue = "1"
const trainingdataNotset = "(not set)"
const trainingdataNewvisitor = "(New Visitor)"
const trainingdataDirect = "(direct)"
const trainingdataReferral = "referral"

func normalizeAnalyticsData(rows []api.AnalyticsDataRow) [][]string {
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

		isReferral := trainingdataFalse
		if row.Medium == trainingdataReferral {
			isReferral = trainingdataTrue
		}

		landingPagePathIsSet := trainingdataTrue
		if row.LandingPagePath == "/" {
			landingPagePathIsSet = trainingdataFalse
		}

		rowValues := []string{
			isNewVisitor,
			fullReferrerIsSet,
			isReferral,
			landingPagePathIsSet,
			strconv.FormatInt(row.Sessions, 10),
			strconv.FormatFloat(row.BounceRate, 'f', -1, 32),
			strconv.FormatFloat(row.PageviewsPerSession, 'f', -1, 32),
			strconv.FormatFloat(row.TimeOnPage, 'f', -1, 32),
			row.Source,
		}

		values = append(values, rowValues)
	}

	return values
}
