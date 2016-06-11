package spamcontrol

import "fmt"

type spamAnalysis interface {
	// GetSpamAnalysis returns a spam analysis report for a given account ID.
	// Returns an error if the report creation failed.
	GetSpamAnalysis(accountID string, numberOfDays int, threshold float64) (AnalysisResult, error)
}

// dynamicSpamAnalysis performs dynamic referrer spam analysis
// on analytics data.
type dynamicSpamAnalysis struct {
	analyticsDataProvider analyticsDataProvider
	spamDetector          SpamDetector
}

// GetSpamAnalysis returns a dynamic referrer spam analysis report for the
// analytics account with the given account ID. Returns an error if the report
// could not be created.
func (spamControl *dynamicSpamAnalysis) GetSpamAnalysis(accountID string, numberOfDays int, threshold float64) (AnalysisResult, error) {

	analyticsData, analyticsDataError := spamControl.analyticsDataProvider.GetAnalyticsData(accountID, numberOfDays)
	if analyticsDataError != nil {
		return AnalysisResult{}, analyticsDataError
	}

	ratedAnalyticsData, spamDetectionError := spamControl.spamDetector.GetSpamRating(analyticsData)
	if spamDetectionError != nil {
		return AnalysisResult{}, spamDetectionError
	}

	// get all spam domains
	spamDomainMap := make(map[string][]SpamDomain)
	for _, row := range ratedAnalyticsData {
		if !row.IsSpam {
			continue
		}

		spamDomainMap[row.Source] = append(spamDomainMap[row.Source], SpamDomain{
			DomainName:      row.Source,
			SpamProbability: row.Probability,
		})
	}

	fmt.Println(spamDomainMap)

	var spamDomains []SpamDomain
	for domainName, domains := range spamDomainMap {

		propability := getAverageProbability(domains)
		if propability < threshold {
			continue
		}

		spamDomains = append(spamDomains, SpamDomain{
			DomainName:      domainName,
			SpamProbability: propability,
		})
	}

	// sort the domains by name
	SortSpamDomainsBy(spamDomainsByName).Sort(spamDomains)

	// assemble a view model
	spamStatusModel := AnalysisResult{
		AccountID:   accountID,
		SpamDomains: spamDomains,
	}

	return spamStatusModel, nil
}
