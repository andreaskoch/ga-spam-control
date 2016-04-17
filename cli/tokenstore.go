package main

import (
	"os"

	"golang.org/x/oauth2"
)

// newTokenStore creates a new token store instance that uses
// the given file path to read and store API tokens.
func newTokenStore(credentialFilePath string) filesystemTokenStore {
	return filesystemTokenStore{
		credentialFilePath,
		tokenSerializer{},
	}
}

// filesystemTokenStore provides functions for reading and persisting
// oAuth tokens in a given file path.
type filesystemTokenStore struct {
	filePath   string
	serializer tokenSerializer
}

// GetToken returns an oAuth token from disc. Returns an
// error if no token was found.
func (store filesystemTokenStore) GetToken() (*oauth2.Token, error) {
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
func (store filesystemTokenStore) SaveToken(token *oauth2.Token) error {
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
