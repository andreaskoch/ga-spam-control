package main

import (
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/oauth2"
)

type tokenSerializer struct{}

func newTokenSerializer() tokenSerializer {
	return tokenSerializer{}
}

func (tokenSerializer) Serialize(writer io.Writer, token oauth2.Token) error {
	bytes, err := json.MarshalIndent(token, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (tokenSerializer) Deserialize(reader io.Reader) (oauth2.Token, error) {
	decoder := json.NewDecoder(reader)
	var token oauth2.Token
	err := decoder.Decode(token)
	return token, err
}

func newTokenStore(credentialFilePath string) tokenStore {
	return tokenStore{credentialFilePath, newTokenSerializer()}
}

type tokenStore struct {
	filePath   string
	serializer tokenSerializer
}

func (store tokenStore) GetToken() (oauth2.Token, error) {
	return oauth2.Token{}, fmt.Errorf("Token not found")
}

func (store tokenStore) SaveToken(token oauth2.Token) error {
	return nil
}
