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

type Table struct {
	ColumnNames []string
	Rows        [][]string
}

func rowsToTable(rows []api.AnalyticsDataRow) Table {

	var values [][]string
	for _, row := range rows {

		isNewVisitor := "0"
		if row.UserType == "New Visitor" {
			isNewVisitor = "1"
		}

		fullReferrerIsSet := "1"
		if row.FullReferrer == "(direct)" {
			fullReferrerIsSet = "0"
		}

		mediumIsSet := "1"
		if row.Medium == "(not set)" {
			mediumIsSet = "0"
		}

		networkDomainIsSet := "1"
		if row.NetworkDomain == "(not set)" {
			networkDomainIsSet = "0"
		}

		networkLocationIsSet := "1"
		if row.NetworkLocation == "(not set)" {
			networkLocationIsSet = "0"
		}

		landingPagePathIsSet := "1"
		if row.LandingPagePath == "/" {
			landingPagePathIsSet = "0"
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
