// Package token implements the filesystem based token-store for the
// Google Analytics API oAuth token.
package token

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/andreaskoch/ga-spam-control/api/apicredentials"
	"github.com/andreaskoch/ga-spam-control/common/fsutil"

	"golang.org/x/oauth2"
)

// NewTokenStore creates a new token store instance that uses
// the given file path to read and store API tokens.
func NewTokenStore(credentialFilePath string) apicredentials.TokenStorer {
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
	folder := filepath.Dir(store.filePath)
	if !fsutil.PathExists(folder) {
		if err := os.MkdirAll(folder, 0700); err != nil {
			return err
		}
	}

	if !fsutil.IsDirectory(folder) {
		return fmt.Errorf("Cannot create folder %q. A file with the same name already exists.", folder)
	}

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
