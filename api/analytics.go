package api

import (
	"fmt"
	"strconv"

	"github.com/andreaskoch/ga-spam-control/api/apiservice"
)

// AnalyticsData is a set of analytics data entries.
type AnalyticsData []AnalyticsDataRow

// AnalyticsDataRow contains a single analytics data record.
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
func toModelAnalyticsData(source apiservice.AnalyticsData) (AnalyticsData, error) {
	analyticsData := AnalyticsData{}

	for _, row := range source.Rows {
		analyticsDataRow := AnalyticsDataRow{}

		for colIndex, col := range source.Cols {

			value := row.Cell[colIndex].Value

			switch col.Label {

			case "ga:userType":
				analyticsDataRow.UserType = value

			case "ga:fullReferrer":
				analyticsDataRow.FullReferrer = value

			case "ga:source":
				analyticsDataRow.Source = value

			case "ga:medium":
				analyticsDataRow.Medium = value

			case "ga:networkDomain":
				analyticsDataRow.NetworkDomain = value

			case "ga:networkLocation":
				analyticsDataRow.NetworkLocation = value

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

			default:
				return AnalyticsData{}, fmt.Errorf("The column %q was not recognized", col.Label)

			}

		}

		analyticsData = append(analyticsData, analyticsDataRow)

	}

	return analyticsData, nil
}
