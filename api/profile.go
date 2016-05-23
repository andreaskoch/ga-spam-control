package api

import "github.com/andreaskoch/ga-spam-control/api/apiservice"

// toModelProfiles converts []apiservice.Profile to []Profile.
func toModelProfiles(sources []apiservice.Profile) []Profile {

	var accounts []Profile
	for _, source := range sources {
		accounts = append(accounts, toModelProfile(source))
	}

	return accounts
}

// toModelProfile converts a apiservice.Profile model into a Profile model.
func toModelProfile(source apiservice.Profile) Profile {
	return Profile{
		ID:                    source.ID,
		Kind:                  source.Kind,
		Link:                  source.SelfLink,
		AccountID:             source.AccountID,
		WebPropertyID:         source.WebPropertyID,
		InternalWebPropertyID: source.InternalWebPropertyID,
		Name:       source.Name,
		Currency:   source.Currency,
		Timezone:   source.Timezone,
		WebsiteURL: source.WebsiteURL,
		Type:       source.Type,
	}
}

// toServiceProfile converts Profile to apiservice.Profile.
func toServiceProfile(source Profile) apiservice.Profile {
	return apiservice.Profile{
		Item: apiservice.Item{
			ID:       source.ID,
			Kind:     source.Kind,
			SelfLink: source.Link,
		},
		AccountID:             source.AccountID,
		WebPropertyID:         source.WebPropertyID,
		InternalWebPropertyID: source.InternalWebPropertyID,
		Name:       source.Name,
		Currency:   source.Currency,
		Timezone:   source.Timezone,
		WebsiteURL: source.WebsiteURL,
		Type:       source.Type,
	}
}

// A Profile contains information about Google Analytics views.
type Profile struct {
	ID   string
	Kind string
	Link string

	AccountID             string
	WebPropertyID         string
	InternalWebPropertyID string
	Name                  string
	Currency              string
	Timezone              string
	WebsiteURL            string
	Type                  string
}
