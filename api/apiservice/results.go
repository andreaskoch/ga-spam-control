package apiservice

type Results struct {
	Kind         string `json:"kind"`
	Username     string `json:"username"`
	TotalResults int    `json:"totalResults"`
	StartIndex   int    `json:"startIndex"`
	ItemsPerPage int    `json:"itemsPerPage"`
}

type Link struct {
	Type string `json:"type"`
	Href string `json:"href"`
}
