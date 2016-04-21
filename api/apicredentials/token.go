package apicredentials

import "golang.org/x/oauth2"

// The TokenStorer interface provides functions for
// reading and persisting oauth2.Token models.
type TokenStorer interface {
	GetToken() (*oauth2.Token, error)
	SaveToken(token *oauth2.Token) error
}
