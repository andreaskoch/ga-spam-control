package main

import (
	"encoding/json"
	"io"
	"os"

	"golang.org/x/oauth2"
)

type tokenSerializer struct{}

func newTokenSerializer() tokenSerializer {
	return tokenSerializer{}
}

func (tokenSerializer) Serialize(writer io.Writer, token *oauth2.Token) error {
	bytes, err := json.MarshalIndent(token, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (tokenSerializer) Deserialize(reader io.Reader) (*oauth2.Token, error) {
	decoder := json.NewDecoder(reader)
	var token *oauth2.Token
	err := decoder.Decode(&token)
	return token, err
}

func newTokenStore(credentialFilePath string) tokenStore {
	return tokenStore{credentialFilePath, newTokenSerializer()}
}

type tokenStore struct {
	filePath   string
	serializer tokenSerializer
}

func (store tokenStore) GetToken() (*oauth2.Token, error) {
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

func (store tokenStore) SaveToken(token *oauth2.Token) error {
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
