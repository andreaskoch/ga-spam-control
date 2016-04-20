package api

import "github.com/andreaskoch/ga-spam-control/api/apiservice"

// toModelProfiles converts []apiservice.Profile to []Profile.
func toModelProfiles(sources []apiservice.Profile) []Profile {

	accounts := make([]Profile, 0)
	for _, source := range sources {
		accounts = append(accounts, toModelProfile(source))
	}

	return accounts
}

// toModelProfile converts a apiservice.Profile model into a Profile model.
func toModelProfile(source apiservice.Profile) Profile {
	return Profile{
		ID:   source.ID,
		Kind: source.Kind,
		Link: source.SelfLink,
		Reference: ProfileReference{
			Name:                  source.Entity.ProfileRef.Name,
			ID:                    source.Entity.ProfileRef.ID,
			Kind:                  source.Entity.ProfileRef.Kind,
			Href:                  source.Entity.ProfileRef.Href,
			AccountID:             source.Entity.ProfileRef.AccountID,
			WebPropertyID:         source.Entity.ProfileRef.WebPropertyID,
			InternalWebPropertyID: source.Entity.ProfileRef.InternalWebPropertyID,
		},
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
		Entity: apiservice.ProfileEntity{
			ProfileRef: apiservice.ProfileReference{
				Name:                  source.Reference.Name,
				ID:                    source.Reference.ID,
				Kind:                  source.Reference.Kind,
				Href:                  source.Reference.Href,
				AccountID:             source.Reference.AccountID,
				WebPropertyID:         source.Reference.WebPropertyID,
				InternalWebPropertyID: source.Reference.InternalWebPropertyID,
			},
		},
	}
}

type Profile struct {
	ID        string
	Kind      string
	Link      string
	Reference ProfileReference
}

type ProfileReference struct {
	ID                    string
	Kind                  string
	Href                  string
	AccountID             string
	WebPropertyID         string
	InternalWebPropertyID string
	Name                  string
}
