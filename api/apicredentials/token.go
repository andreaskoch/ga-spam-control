// Package apicredentials contains the TokenStorer interface that
// is used as the oAuth token provider for the Google Analytics API.
package apicredentials

import "golang.org/x/oauth2"

// The TokenStorer interface provides functions for
// reading and persisting oauth2.Token models.
type TokenStorer interface {
	GetToken() (*oauth2.Token, error)
	SaveToken(token *oauth2.Token) error
}
