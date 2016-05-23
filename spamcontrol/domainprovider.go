package spamcontrol

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

// A SpamDomainProvider provides domain names of
// referrer spam providers.
type SpamDomainProvider interface {
	// GetSpamDomains returns a list of referrer spam domain names.
	GetSpamDomains() ([]string, error)
}

// The staticSpamDomains provides a static list
// of referrer spam domain names from a remote URL.
type staticSpamDomains struct {
	domainListURL string
}

// GetSpamDomains returns a list of referrer spam domain names for a
// remote text file that contains domain names that are considered spam.
func (provider *staticSpamDomains) GetSpamDomains() ([]string, error) {

	// request the domain names from the remote source
	response, requestError := http.Get(provider.domainListURL)
	if requestError != nil {
		return nil, fmt.Errorf("Failed to get URL %q: %s", provider.domainListURL, requestError.Error())
	}

	// read the domain names line-by-line
	var domains []string
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())

		// ignore empty lines
		if domain == "" {
			continue
		}

		domains = append(domains, domain)
	}

	return domains, nil
}
