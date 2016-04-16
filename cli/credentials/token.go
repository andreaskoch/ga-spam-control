package credentials

import (
	"encoding/json"
	"io"
	"os"

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

// NewTokenStore creates a new token store instance that uses
// the given file path to read and store API tokens.
func NewTokenStore(credentialFilePath string) TokenStore {
	return TokenStore{
		credentialFilePath,
		tokenSerializer{},
	}
}

// TokenStore provides functions for reading and persisting
// oAuth tokens in a given file path.
type TokenStore struct {
	filePath   string
	serializer tokenSerializer
}

// GetToken returns an oAuth token from disc. Returns an
// error if no token was found.
func (store TokenStore) GetToken() (*oauth2.Token, error) {
	file, readErr := os.Open(store.filePath)
	if readErr != nil {
		return nil, readErr
	}

	token, deserializeErr := store.serializer.Deserialize(file)
	if deserializeErr != nil {
		return nil, deserializeErr
	}

	return token, nil
}

// SaveToken stores the given oAuth token to disc. Returns an
// error if the save failed.
func (store TokenStore) SaveToken(token *oauth2.Token) error {
	file, fileErr := os.OpenFile(store.filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if fileErr != nil {
		return fileErr
	}

	serializeErr := store.serializer.Serialize(file, token)
	if serializeErr != nil {
		return serializeErr
	}

	return nil
}
