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
	spamDomainProvider := &dummyDomainProvider{domains}
	filterFactory := spamFilterFactory{
		domainProvider:       spamDomainProvider,
		filterNameProvider:   &spamFilterNameProvider{"ga-spam-control"},
		filterValueMaxLength: 255,
	}

	// act
	filters, _ := filterFactory.GetNewFilters()

	// assert
	if len(filters) > 0 {
		t.Fail()
	}
}

func Test_GetNewFilters_ValidDomains_FilterIsReturned(t *testing.T) {
	// arrange
	domains := []string{"referrer-spam.com", "referrer-spam.co.uk"}
	spamDomainProvider := &dummyDomainProvider{domains}
	filterFactory := spamFilterFactory{
		domainProvider:       spamDomainProvider,
		filterNameProvider:   &spamFilterNameProvider{"ga-spam-control"},
		filterValueMaxLength: 255,
	}

	// act
	filters, _ := filterFactory.GetNewFilters()

	// assert
	if len(filters) != 1 {
		t.Fail()
	}
}

func Test_GetNewFilters_ValidDomains_FilterExpressionValueIsCorrect(t *testing.T) {
	// arrange
	domains := []string{"referrer-spam.com", "referrer-spam.co.uk"}
	spamDomainProvider := &dummyDomainProvider{domains}
	filterFactory := spamFilterFactory{
		domainProvider:       spamDomainProvider,
		filterNameProvider:   &spamFilterNameProvider{"ga-spam-control"},
		filterValueMaxLength: 255,
	}

	// act
	filters, _ := filterFactory.GetNewFilters()

	// assert
	filter := filters[0]
	if filter.ExcludeDetails.ExpressionValue != `referrer-spam\.com|referrer-spam\.co\.uk` {
		t.Fail()
		t.Logf("The expression value is invalid: %#v", filter)
	}
}

func Test_GetNewFilters_ValidDomains_FilterNameIsCorrect(t *testing.T) {
	// arrange
	domains := []string{"referrer-spam.com", "referrer-spam.co.uk"}
	spamDomainProvider := &dummyDomainProvider{domains}
	filterFactory := spamFilterFactory{
		domainProvider:       spamDomainProvider,
		filterNameProvider:   &spamFilterNameProvider{"ga-spam-control"},
		filterValueMaxLength: 255,
	}

	// act
	filters, _ := filterFactory.GetNewFilters()

	// assert
	filter := filters[0]
	if filter.Name != `ga-spam-control 01` {
		t.Fail()
		t.Logf("The filter name is invalid: %#v", filter.Name)
	}
}

func Test_GetNewFilters_ManyDomains_ThreeFiltersAreReturned(t *testing.T) {
	// arrange
	domains := []string{
		"0n-line.tv",
		"100dollars-seo.com",
		"12masterov.com",
		"1pamm.ru",
		"4webmasters.org",
	}
	spamDomainProvider := &dummyDomainProvider{domains}
	filterFactory := spamFilterFactory{
		domainProvider:       spamDomainProvider,
		filterNameProvider:   &spamFilterNameProvider{"ga-spam-control"},
		filterValueMaxLength: 35,
	}

	// act
	filters, _ := filterFactory.GetNewFilters()

	// assert
	if len(filters) != 3 {
		t.Fail()
		t.Logf("GetNewFilters did not return 3 segments: %#v", filters)
	}
}

func Test_GetNewFilters_TooLongDomainName_ErrorIsReturned(t *testing.T) {
	// arrange
	domains := []string{
		"1234567890-too-long-to-fit.com",
	}
	spamDomainProvider := &dummyDomainProvider{domains}
	filterFactory := spamFilterFactory{
		domainProvider:       spamDomainProvider,
		filterNameProvider:   &spamFilterNameProvider{"ga-spam-control"},
		filterValueMaxLength: 10,
	}

	// act
	_, err := filterFactory.GetNewFilters()

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
	segments, _ := getExpressionSegments(domainNames, 10)

	// assert
	if len(segments) > 0 {
		t.Fail()
		t.Logf("getExpressionSegments should not return any segments: %#v", segments)
	}
}

func Test_getExpressionSegments_DomainNameIsEscaped(t *testing.T) {
	// arrange
	domainNames := []string{"0a-zine.tv"}

	// act
	segments, _ := getExpressionSegments(domainNames, 15)

	// assert
	if len(segments) != 1 || segments[0] != `0a-zine\.tv` {
		t.Fail()
		t.Logf("getExpressionSegments did not escape the domain name: %#v", segments)
	}
}

func Test_getExpressionSegments_3SegmentsAreCreated(t *testing.T) {
	// arrange
	domainNames := []string{
		"0n-line.tv",
		"100dollars-seo.com",
		"12masterov.com",
		"1pamm.ru",
		"4webmasters.org",
	}

	maxSegmentSize := 35

	// act
	segments, _ := getExpressionSegments(domainNames, maxSegmentSize)

	// assert
	if len(segments) != 3 {
		t.Fail()
		t.Logf("getExpressionSegments did not return 3 segments: %#v", segments)
	}
}

func Test_getExpressionSegments_SegmentsDontFit_ErrorIsReturned(t *testing.T) {
	// arrange
	domainNames := []string{
		"123.com",
		"123456789.com",
		"123456789.co.uk",
	}

	maxSegmentSize := 8

	// act
	_, err := getExpressionSegments(domainNames, maxSegmentSize)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("getExpressionSegments did not return an error.")
	}
}

func Test_getExpressionSegments_SegmentsAreConnectedWithPipes(t *testing.T) {
	// arrange
	domainNames := []string{
		"0n-line.tv",
		"100dollars-seo.com",
	}

	maxSegmentSize := 35

	// act
	segments, _ := getExpressionSegments(domainNames, maxSegmentSize)

	// assert
	if len(segments) != 1 || segments[0] != `0n-line\.tv|100dollars-seo\.com` {
		t.Fail()
		t.Logf("getExpressionSegments did not escape properly: %s", segments[0])
	}
}
