package main

import (
	"encoding/json"
	"io"

	"golang.org/x/oauth2"
)

// tokenSerializer contains functions for serializing an
// de-serializing oauth2.Token models.
type tokenSerializer struct{}

// Serialize the given oauth2.Token to the given writer.
// Returns an error if the serialization failed.
func (tokenSerializer) Serialize(writer io.Writer, token *oauth2.Token) error {
	bytes, err := json.MarshalIndent(token, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

// Deserialize the oauth2.Token model from the given io.Reader.
func (tokenSerializer) Deserialize(reader io.Reader) (*oauth2.Token, error) {
	decoder := json.NewDecoder(reader)
	var token *oauth2.Token
	err := decoder.Decode(&token)
	return token, err
}
