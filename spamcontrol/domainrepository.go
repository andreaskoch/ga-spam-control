package spamcontrol

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/andreaskoch/ga-spam-control/common/fsutil"
)

// NewCommunitySpamDomainRepository creates a new CommunitySpamDomainRepository instance.
func NewCommunitySpamDomainRepository(filePath string, providers []SpamDomainProvider) *CommunitySpamDomainRepository {
	return &CommunitySpamDomainRepository{
		LocalSpamDomainRepository{
			filePath:  filePath,
			providers: providers,
		},
	}
}

// NewPrivateSpamDomainRepository creates a new PrivateSpamDomainRepository instance.
func NewPrivateSpamDomainRepository(filePath string, provider SpamDomainProvider) *PrivateSpamDomainRepository {
	return &PrivateSpamDomainRepository{
		LocalSpamDomainRepository{
			filePath:  filePath,
			providers: []SpamDomainProvider{provider},
		},
	}
}

// A CommunitySpamDomainRepository stores the referrer spam domains names
// that are provided by the community locally and provides functions for
// keeping this file up-to-date.
type CommunitySpamDomainRepository struct {
	LocalSpamDomainRepository
}

// GetSpamDomains returns a list of referrer spam domains from disc.
func (repository *CommunitySpamDomainRepository) GetSpamDomains() ([]string, error) {
	return repository.LocalSpamDomainRepository.GetSpamDomains()
}

// UpdateSpamDomains stores the latest list of spam domains to disc.
func (repository *CommunitySpamDomainRepository) UpdateSpamDomains() (unchanged, added, removed []string, err error) {
	return repository.LocalSpamDomainRepository.UpdateSpamDomains()
}

// A PrivateSpamDomainRepository provides functions for saving a personal
// list of referrer spam domain names locally.
type PrivateSpamDomainRepository struct {
	LocalSpamDomainRepository
}

// GetSpamDomains returns a list of referrer spam domains from disc.
func (repository *PrivateSpamDomainRepository) GetSpamDomains() ([]string, error) {
	return repository.LocalSpamDomainRepository.GetSpamDomains()
}

// AddDomains add the given spam domain names to the repository.
func (repository *PrivateSpamDomainRepository) AddDomains(newDomains []string) error {
	existingDomains, existingDomainsError := repository.LocalSpamDomainRepository.GetSpamDomains()
	if existingDomainsError != nil {
		return existingDomainsError
	}

	mergedDomainNames := append(existingDomains, newDomains...)
	sort.Strings(mergedDomainNames)
	mergedDomainNames = unique(mergedDomainNames)

	writeError := repository.LocalSpamDomainRepository.Write(mergedDomainNames)
	if writeError != nil {
		return writeError
	}

	return nil
}

// A LocalSpamDomainRepository provides functions
// for retrieving and storing spam domain names on disc.
type LocalSpamDomainRepository struct {
	filePath  string
	providers []SpamDomainProvider
}

// GetSpamDomains returns a list of referrer spam domains from disc.
func (repository *LocalSpamDomainRepository) GetSpamDomains() ([]string, error) {
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
func (repository *LocalSpamDomainRepository) UpdateSpamDomains() (unchanged, added, removed []string, err error) {

	// get the existing domains for comparision
	existingDomains, _ := getDomainsFromFile(repository.filePath)

	// get new spam domains
	newDomains, domainError := repository.getLatestSpamDomains()
	if domainError != nil {
		return nil, nil, nil, domainError
	}

	// write the new list to disc
	writeError := repository.Write(newDomains)
	if writeError != nil {
		return nil, nil, nil, writeError
	}

	unchanged, added, removed = diff(existingDomains, newDomains)
	return
}

// Write stores the given domain names in the repository.
func (repository *LocalSpamDomainRepository) Write(domains []string) error {
	writeError := writeDomainsToFile(repository.filePath, domains)
	if writeError != nil {
		return writeError
	}

	return nil
}

func (repository *LocalSpamDomainRepository) getLatestSpamDomains() ([]string, error) {

	var fullDomainList []string
	for _, provider := range repository.providers {
		domains, err := provider.GetSpamDomains()
		if err != nil {
			return nil, err
		}

		fullDomainList = append(fullDomainList, domains...)
	}

	fullDomainList = unique(fullDomainList)

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

	// make sure the folder exists
	folder := filepath.Dir(filePath)
	if !fsutil.PathExists(folder) {
		if err := os.MkdirAll(folder, 0700); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0660)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, domain := range domains {
		fmt.Fprintf(writer, "%s"+NewLineSequence, domain)
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
