package apiservice

type Results struct {
	Kind         string `json:"kind"`
	Username     string `json:"username"`
	TotalResults int    `json:"totalResults"`
	StartIndex   int    `json:"startIndex"`
	ItemsPerPage int    `json:"itemsPerPage"`
}

type Item struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	SelfLink string `json:"selfLink"`
}

type Link struct {
	Type string `json:"type"`
	Href string `json:"href"`
}
