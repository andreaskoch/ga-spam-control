package apiservice

import (
	"encoding/json"
	"io"
)

type profileResultsSerializer struct{}

func (profileResultsSerializer) Serialize(writer io.Writer, profileResults *ProfileResults) error {
	bytes, err := json.MarshalIndent(profileResults, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (profileResultsSerializer) Deserialize(reader io.Reader) (*ProfileResults, error) {
	decoder := json.NewDecoder(reader)
	var profileResults *ProfileResults
	err := decoder.Decode(&profileResults)
	return profileResults, err
}

type profileSerializer struct{}

func (profileSerializer) Serialize(writer io.Writer, profile *Profile) error {
	bytes, err := json.MarshalIndent(profile, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (profileSerializer) Deserialize(reader io.Reader) (*Profile, error) {
	decoder := json.NewDecoder(reader)
	var profile *Profile
	err := decoder.Decode(&profile)
	return profile, err
}

type ProfileResults struct {
	Results
	Items []Profile `json:"items"`
}

type Profile struct {
	Item
	AccountID             string `json:"accountId"`
	WebPropertyID         string `json:"webPropertyId"`
	InternalWebPropertyID string `json:"internalWebPropertyId"`
	Name                  string `json:"name"`
	Currency              string `json:"currency"`
	Timezone              string `json:"timezone"`
	WebsiteURL            string `json:"websiteUrl"`
	Type                  string `json:"type"`
}
