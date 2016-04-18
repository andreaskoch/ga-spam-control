package spamcontrol

// A spamDomainProvider provides domain names of
// referer spam providers.
type spamDomainProvider interface {
	// GetSpamDomains returns a list of referer spam domain names.
	GetSpamDomains() ([]string, error)
}

// The remoteSpamDomainProvider fetches the list
// of referer spam domain names from a remote URL.
type remoteSpamDomainProvider struct {
	domainListUrl string
}

// GetSpamDomains returns a list of referer spam domain names.
func (spamDomainProvider *remoteSpamDomainProvider) GetSpamDomains() ([]string, error) {
	return []string{"0n-line.tv", "100dollars-seo.com"}, nil
}
