package api

import (
	"strconv"

	"github.com/andreaskoch/ga-spam-control/api/apiservice"
)

type AnalyticsData struct {
	Rows []AnalyticsDataRow
}

type AnalyticsDataRow struct {
	// Dimensions
	UserType        string
	FullReferrer    string
	Source          string
	Medium          string
	NetworkDomain   string
	NetworkLocation string
	LandingPagePath string

	// Metrics
	Sessions            int64
	BounceRate          float64
	PageviewsPerSession float64
	TimeOnPage          float64
}

// toModelAnalyticsData converts a apiservice.AnalyticsData model into a AnalyticsData model.
func toModelAnalyticsData(source apiservice.AnalyticsData) AnalyticsData {
	analyticsData := AnalyticsData{}

	for _, row := range source.Rows {
		analyticsDataRow := AnalyticsDataRow{}

		for colIndex, col := range source.Cols {

			value := row.Cell[colIndex].Value
			if value == "" {
				value = "(not set)"
			}

			switch col.Label {

			case "ga:userType":
				analyticsDataRow.UserType = value

			case "ga:fullReferrer":
				analyticsDataRow.FullReferrer = value

			case "ga:source":
				analyticsDataRow.Source = value

			case "ga:medim":
				analyticsDataRow.Medium = value

			case "ga:networkDomain":
				analyticsDataRow.NetworkDomain = value

			case "ga:landingPagePath":
				analyticsDataRow.LandingPagePath = value

			case "ga:sessions":
				intValue, _ := strconv.ParseInt(value, 10, 64)
				analyticsDataRow.Sessions = intValue

			case "ga:bounceRate":
				floatValue, _ := strconv.ParseFloat(value, 64)
				analyticsDataRow.BounceRate = floatValue

			case "ga:pageviewsPerSession":
				floatValue, _ := strconv.ParseFloat(value, 64)
				analyticsDataRow.PageviewsPerSession = floatValue

			case "ga:timeOnPage":
				floatValue, _ := strconv.ParseFloat(value, 64)
				analyticsDataRow.TimeOnPage = floatValue

			}

			analyticsData.Rows = append(analyticsData.Rows, analyticsDataRow)

		}

	}

	return analyticsData
}
