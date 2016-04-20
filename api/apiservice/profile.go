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

type ProfileDetail struct {
	Kind            string `json:"kind"`
	Field           string `json:"field"`
	MatchType       string `json:"matchType"`
	ExpressionValue string `json:"expressionValue"`
	CaseSensitive   bool   `json:"caseSensitive"`
}

type Profile struct {
	Item
	Entity ProfileEntity `json:"entity"`
}

type ProfileEntity struct {
	ProfileRef ProfileReference `json:"profileRef"`
}

type ProfileReference struct {
	ID                    string `json:"id"`
	Kind                  string `json:"kind"`
	Href                  string `json:"href"`
	AccountID             string `json:"accountId"`
	WebPropertyID         string `json:"webPropertyId"`
	InternalWebPropertyID string `json:"internalWebPropertyId"`
	Name                  string `json:"name"`
}
