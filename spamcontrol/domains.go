package spamcontrol

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

// A spamDomainProvider provides domain names of
// referrer spam providers.
type spamDomainProvider interface {
	// GetSpamDomains returns a list of referrer spam domain names.
	GetSpamDomains() ([]string, error)
}

// The remoteSpamDomainProvider fetches the list
// of referrer spam domain names from a remote URL.
type remoteSpamDomainProvider struct {
	domainListURL string
}

// GetSpamDomains returns a list of referrer spam domain names.
func (spamDomainProvider *remoteSpamDomainProvider) GetSpamDomains() ([]string, error) {

	// request the domain names from the remote source
	response, requestError := http.Get(spamDomainProvider.domainListURL)
	if requestError != nil {
		return nil, fmt.Errorf("Failed to get URL %q: %s", spamDomainProvider.domainListURL, requestError.Error())
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
