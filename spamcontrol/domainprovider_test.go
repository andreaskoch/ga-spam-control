package spamcontrol

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetSpamDomains_200OK_ResponseContainsOneLine_SingleDomainIsReturned(t *testing.T) {
	// arrange
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "example.com")
	}))
	defer testServer.Close()

	domainProvider := staticSpamDomains{testServer.URL}

	// act
	domains, err := domainProvider.GetSpamDomains()

	// assert
	if err != nil {
		t.Fail()
		t.Logf("GetSpamDomains() returned an error: %s", err.Error())
	}

	if len(domains) != 1 || domains[0] != "example.com" {
		t.Fail()
		t.Logf("GetSpamDomains() did not return a domain")
	}
}

func Test_GetSpamDomains_200OK_ResponseContainsTwoLines_TwoDomainsAreReturned(t *testing.T) {
	// arrange
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "example1.com")
		fmt.Fprintln(w, "example2.com")
	}))
	defer testServer.Close()

	domainProvider := staticSpamDomains{testServer.URL}

	// act
	domains, err := domainProvider.GetSpamDomains()

	// assert
	if err != nil {
		t.Fail()
		t.Logf("GetSpamDomains() returned an error: %s", err.Error())
	}

	if len(domains) != 2 || domains[0] != "example1.com" || domains[1] != "example2.com" {
		t.Fail()
		t.Logf("GetSpamDomains() did not return the expected domains: %v", domains)
	}
}

func Test_GetSpamDomains_200OK_ResponseContainsDiverseText_DomainNamesAreNotModified(t *testing.T) {
	// arrange
	testRecords := []struct {
		Input          string
		ExpectedOutput string
	}{
		{"ExamPLE.com", "ExamPLE.com"},
		{"пример.ru", "пример.ru"},
		{"例.cn", "例.cn"},
		{"häst.se", "häst.se"},
		{"fråga.se", "fråga.se"},
		{"Påskeøl.dk", "Påskeøl.dk"},
	}

	for _, testRecord := range testRecords {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, testRecord.Input)
		}))
		defer testServer.Close()

		domainProvider := staticSpamDomains{testServer.URL}

		// act
		domains, err := domainProvider.GetSpamDomains()

		// assert
		if err != nil {
			t.Fail()
			t.Logf("GetSpamDomains() returned an error: %s", err.Error())
		}

		if len(domains) != 1 || domains[0] != testRecord.ExpectedOutput {
			t.Fail()
			t.Logf("GetSpamDomains() did not returned the expected result %q", testRecord.ExpectedOutput)
		}
	}
}

func Test_GetSpamDomains_200OK_ResponseContainsDomainsWithWhitespace_DomainnamesAreTrimmed(t *testing.T) {
	// arrange
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "example1.com")
		fmt.Fprintln(w, "example2.com ")
		fmt.Fprintln(w, " example3.com")
	}))
	defer testServer.Close()

	domainProvider := staticSpamDomains{testServer.URL}

	// act
	domains, err := domainProvider.GetSpamDomains()

	// assert
	if err != nil {
		t.Fail()
		t.Logf("GetSpamDomains() returned an error: %s", err.Error())
	}

	if len(domains) != 3 || domains[0] != "example1.com" || domains[1] != "example2.com" || domains[2] != "example3.com" {
		t.Fail()
		t.Logf("GetSpamDomains() did not return the expected domains: %v", domains)
	}
}

func Test_GetSpamDomains_200OK_ResponseContainsEmptyLines_EmptyLinesAreOmitted(t *testing.T) {
	// arrange
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "example1.com")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "example3.com")
	}))
	defer testServer.Close()

	domainProvider := staticSpamDomains{testServer.URL}

	// act
	domains, err := domainProvider.GetSpamDomains()

	// assert
	if err != nil {
		t.Fail()
		t.Logf("GetSpamDomains() returned an error: %s", err.Error())
	}

	if len(domains) != 2 || domains[0] != "example1.com" || domains[1] != "example3.com" {
		t.Fail()
		t.Logf("GetSpamDomains() did not return the expected domains: %v", domains)
	}
}

func Test_GetSpamDomains_404Error_ResponseContains404ErrorPage_ErrorIsReturned(t *testing.T) {
	// arrange
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)

		fmt.Fprintln(w, "Error 404 la di da")
	}))
	defer testServer.Close()

	domainProvider := staticSpamDomains{testServer.URL}

	// act
	domains, err := domainProvider.GetSpamDomains()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetSpamDomains() did not return an error.")
	}

	if len(domains) != 0 {
		t.Fail()
		t.Logf("GetSpamDomains() did return domains even though the service returned a 404 error: %v", domains)
	}
}

func Test_GetSpamDomains_500Error_ResponseContainsErrorPage_ErrorIsReturned(t *testing.T) {
	// arrange
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "500 Internal Server Error", 500)

		fmt.Fprintln(w, "Error 500 la di da")
	}))
	defer testServer.Close()

	domainProvider := staticSpamDomains{testServer.URL}

	// act
	domains, err := domainProvider.GetSpamDomains()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetSpamDomains() did not return an error.")
	}

	if len(domains) != 0 {
		t.Fail()
		t.Logf("GetSpamDomains() did return domains even though the service returned an error: %v", domains)
	}
}
