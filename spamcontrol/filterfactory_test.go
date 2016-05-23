package spamcontrol

import "testing"

type dummyDomainProvider struct {
	domainNames []string
}

func (domainProvider dummyDomainProvider) GetSpamDomains() ([]string, error) {
	return domainProvider.domainNames, nil
}

func Test_GetNewFilters_NoDomains_NoFilters(t *testing.T) {
	// arrange
	domains := []string{}
	filterFactory := googleAnalyticsFilterFactory{
		filterNameProvider:   &spamFilterNameProvider{"ga-spam-control"},
		filterValueMaxLength: 255,
	}

	// act
	filters, _ := filterFactory.GetNewFilters(domains)

	// assert
	if len(filters) > 0 {
		t.Fail()
	}
}

func Test_GetNewFilters_ValidDomains_OneFilterIsReturned(t *testing.T) {
	// arrange
	domains := []string{"referrer-spam.com", "referrer-spam.co.uk"}
	filterFactory := googleAnalyticsFilterFactory{
		filterNameProvider:   &spamFilterNameProvider{"ga-spam-control"},
		filterValueMaxLength: 255,
	}

	// act
	filters, _ := filterFactory.GetNewFilters(domains)

	// assert
	if len(filters) != 1 {
		t.Fail()
	}
}

func Test_GetNewFilters_ValidDomains_FilterExpressionValueIsCorrect(t *testing.T) {
	// arrange
	domains := []string{"referrer-spam.com", "referrer-spam.co.uk"}
	filterFactory := googleAnalyticsFilterFactory{
		filterNameProvider:   &spamFilterNameProvider{"ga-spam-control"},
		filterValueMaxLength: 255,
	}

	// act
	filters, _ := filterFactory.GetNewFilters(domains)

	// assert
	// note the sorting
	filter := filters[0]
	t.Logf("%q", filter.ExcludeDetails.ExpressionValue)
	if filter.ExcludeDetails.ExpressionValue != "referrer-spam\\.co\\.uk|referrer-spam\\.com" {
		t.Fail()
		t.Logf("The expression value is invalid: %#v", filter)
	}
}

func Test_GetNewFilters_ValidDomains_FilterNameIsCorrect(t *testing.T) {
	// arrange
	domains := []string{"referrer-spam.com", "referrer-spam.co.uk"}
	filterFactory := googleAnalyticsFilterFactory{
		filterNameProvider:   &spamFilterNameProvider{"ga-spam-control"},
		filterValueMaxLength: 255,
	}

	// act
	filters, _ := filterFactory.GetNewFilters(domains)

	// assert
	filter := filters[0]
	if filter.Name != `ga-spam-control Segment R #001` {
		t.Fail()
		t.Logf("The filter name is invalid: %#v", filter.Name)
	}
}

func Test_GetNewFilters_ManyDomains_OneSegment_ThreeFiltersAreReturned(t *testing.T) {
	// arrange
	domains := []string{
		"1n-line.tv",
		"100dollars-seo.com",
		"12masterov.com",
		"1pamm.ru",
		"1webmasters.org",
	}
	filterFactory := googleAnalyticsFilterFactory{
		filterNameProvider:   &spamFilterNameProvider{"ga-spam-control"},
		filterValueMaxLength: 35,
	}

	// act
	filters, _ := filterFactory.GetNewFilters(domains)

	// assert
	if len(filters) != 3 {
		t.Fail()
		t.Logf("GetNewFilters did not return 3 but %d segments", len(filters))
	}
}

func Test_GetNewFilters_ManyDomains_ThreeSegments_FourFiltersAreReturned(t *testing.T) {
	// arrange
	domains := []string{
		"0n-line.tv",
		"100dollars-seo.com",
		"12masterov.com",
		"1pamm.ru",
		"4webmasters.org",
	}
	filterFactory := googleAnalyticsFilterFactory{
		filterNameProvider:   &spamFilterNameProvider{"ga-spam-control"},
		filterValueMaxLength: 35,
	}

	// act
	filters, _ := filterFactory.GetNewFilters(domains)

	// assert
	if len(filters) != 4 {
		t.Fail()
		t.Logf("GetNewFilters did not return 4 but %d segments", len(filters))
	}
}

func Test_GetNewFilters_TooLongDomainName_ErrorIsReturned(t *testing.T) {
	// arrange
	domains := []string{
		"1234567890-too-long-to-fit.com",
	}
	filterFactory := googleAnalyticsFilterFactory{
		filterNameProvider:   &spamFilterNameProvider{"ga-spam-control"},
		filterValueMaxLength: 10,
	}

	// act
	_, err := filterFactory.GetNewFilters(domains)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetNewFilters did not return an error.")
	}
}

func Test_getExpressionSegments_NoDomainNames_NoSegments(t *testing.T) {
	// arrange
	domainNames := []string{}

	// act
	segments, _ := getFilterShards(domainNames, 10)

	// assert
	if len(segments) > 0 {
		t.Fail()
		t.Logf("getExpressionSegments should not return any segments: %#v", segments)
	}
}
