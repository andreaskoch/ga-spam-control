package apiservice

import "time"

type Results struct {
	Kind         string `json:"kind"`
	Username     string `json:"username"`
	TotalResults int    `json:"totalResults"`
	StartIndex   int    `json:"startIndex"`
	ItemsPerPage int    `json:"itemsPerPage"`
}

type Item struct {
	ID       string    `json:"id"`
	Kind     string    `json:"kind"`
	SelfLink string    `json:"selfLink"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type Link struct {
	Type string `json:"type"`
	Href string `json:"href"`
}
