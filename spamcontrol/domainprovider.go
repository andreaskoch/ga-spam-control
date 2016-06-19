package spamcontrol

import (
	"bufio"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// A SpamDomainProvider provides domain names of
// referrer spam providers.
type SpamDomainProvider interface {
	// GetSpamDomains returns a list of referrer spam domain names.
	GetSpamDomains() ([]string, error)
}

// NewAggregateProvider returns a new aggregate providers for the given list of SpamDomainProviders.
func NewAggregateProvider(providers []SpamDomainProvider) SpamDomainProvider {
	return AggregateProvider(providers)
}

// NewRemoteSpamDomainProvider creates a new remote spam domain provider for the given URL.
func NewRemoteSpamDomainProvider(url string) SpamDomainProvider {
	return &remoteSpamDomains{url}
}

// NewLocalSpamDomainProvider creates a new remote spam domain provider for the file path.
func NewLocalSpamDomainProvider(filePath string) SpamDomainProvider {
	return &localSpamDomains{filePath}
}

// An AggregateProvider combines multiple SpamDomainProviders.
type AggregateProvider []SpamDomainProvider

// GetSpamDomains returns a combined list of referrer spam domains from all providers.
func (aggregateProvider AggregateProvider) GetSpamDomains() ([]string, error) {

	var domains []string

	for _, provider := range aggregateProvider {
		repositoryDomains, err := provider.GetSpamDomains()
		if err != nil {
			return nil, err
		}

		domains = append(domains, repositoryDomains...)
	}

	sort.Strings(domains)
	domains = unique(domains)

	return domains, nil
}

// The remoteSpamDomains provides a static list
// of referrer spam domain names from a remote URL.
type remoteSpamDomains struct {
	domainListURL string
}

// GetSpamDomains returns a list of referrer spam domain names from a remote text file.
func (provider *remoteSpamDomains) GetSpamDomains() ([]string, error) {

	// request the domain names from the remote source
	response, requestError := http.Get(provider.domainListURL)
	if requestError != nil {
		return nil, fmt.Errorf("Failed to get URL %q: %s", provider.domainListURL, requestError.Error())

	}

	// check the HTTP status code
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get URL %q. Received HTTP status code %d.", provider.domainListURL, response.StatusCode)
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

// The localSpamDomains provides a static list
// of referrer spam domain names from a local file.
type localSpamDomains struct {
	filePath string
}

// GetSpamDomains returns a list of referrer spam domain names from a local file.
func (provider *localSpamDomains) GetSpamDomains() ([]string, error) {
	return nil, nil
}
