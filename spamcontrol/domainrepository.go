package spamcontrol

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// The SpamDomainRepository interface provides functions for
// retrieving and storing spam domain names.
type SpamDomainRepository interface {

	// GetSpamDomains returns a list of referrer spam domains.
	GetSpamDomains() ([]string, error)

	// UpdateSpamDomains stores the latest list of spam domains.
	// Returns the additions, removals or an error if the update failed.
	UpdateSpamDomains() (unchanged, added, removed []string, err error)
}

// NewSpamDomainRepository creates a new FilesystemSpamDomainRepository instance.
func NewSpamDomainRepository(filePath string, domainProviderFactory *SpamDomainProviderFactory) SpamDomainRepository {
	return &FilesystemSpamDomainRepository{
		filePath:              filePath,
		domainProviderFactory: domainProviderFactory,
	}
}

// A FilesystemSpamDomainRepository provides functions
// for retrieving and storing spam domain names on disc.
type FilesystemSpamDomainRepository struct {
	filePath string

	domainProviderFactory *SpamDomainProviderFactory
}

// GetSpamDomains returns a list of referrer spam domains from disc.
func (repository *FilesystemSpamDomainRepository) GetSpamDomains() ([]string, error) {
	domains, err := getDomainsFromFile(repository.filePath)
	if err != nil {

		_, additions, _, updateError := repository.UpdateSpamDomains()
		if updateError != nil {
			return nil, updateError
		}

		domains = additions
	}

	return domains, nil
}

// UpdateSpamDomains stores the latest list of spam domains to disc.
func (repository *FilesystemSpamDomainRepository) UpdateSpamDomains() (unchanged, added, removed []string, err error) {

	// get the existing domains for comparision
	existingDomains, _ := getDomainsFromFile(repository.filePath)

	// get new spam domains
	newDomains, domainError := repository.getLatestSpamDomains()
	if domainError != nil {
		return nil, nil, nil, domainError
	}

	// write the new list to disc
	writeError := writeDomainsToFile(repository.filePath, newDomains)
	if writeError != nil {
		return nil, nil, nil, writeError
	}

	unchanged, added, removed = diff(existingDomains, newDomains)
	return
}

func (repository *FilesystemSpamDomainRepository) getLatestSpamDomains() ([]string, error) {
	providers, err := repository.domainProviderFactory.GetSpamDomainProviders()
	if err != nil {
		return nil, err
	}

	var fullDomainList []string
	for _, provider := range providers {
		domains, err := provider.GetSpamDomains()
		if err != nil {
			return nil, err
		}

		fullDomainList = append(fullDomainList, domains...)
	}

	fullDomainList = removeDuplicatesFromList(fullDomainList)

	sort.Strings(fullDomainList)

	return fullDomainList, nil
}

// diff calculates the additionas and removals from slice 1 and 2.
func diff(slice1, slice2 []string) (unchanged, added, removed []string) {
	index1 := make(map[string]int)
	for _, member := range slice1 {
		index1[member]++
	}

	index2 := make(map[string]int)
	for _, member := range slice2 {
		index2[member]++
	}

	// detect removals
	for _, member := range slice1 {
		if _, isMemberOfSlice2 := index2[member]; !isMemberOfSlice2 {
			removed = append(removed, member)
		}
	}

	// detect additions
	for _, member := range slice2 {
		if _, isMemberOfSlice1 := index1[member]; !isMemberOfSlice1 {
			added = append(added, member)
		} else {
			unchanged = append(unchanged, member)
		}
	}

	return
}

// writeDomainsToFile writes a list of domain names to the specified file.
// Returns an error if the write failed.
func writeDomainsToFile(filePath string, domains []string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0660)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, domain := range domains {
		fmt.Fprintf(writer, "%s\n", domain)
	}

	return nil
}

// getDomainsFromFile reads referrer spam domains in the given
// filePath. Returns an error if the file was not found, could
// not be read or is empty.
func getDomainsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	// read the file line-by-line
	var domains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content := strings.TrimSpace(scanner.Text())
		if content == "" {
			continue
		}

		domains = append(domains, content)
	}

	// handle errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(domains) == 0 {
		return nil, fmt.Errorf("No domains found in %q", filePath)
	}

	return domains, nil
}
