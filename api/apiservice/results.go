package apiservice

// Results contains common attributes of Google Analtics API results.
type Results struct {
	Kind         string `json:"kind"`
	Username     string `json:"username"`
	TotalResults int    `json:"totalResults"`
	StartIndex   int    `json:"startIndex"`
	ItemsPerPage int    `json:"itemsPerPage"`
}

// Item contains generic attributes of Google Analytics API response objects.
type Item struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	SelfLink string `json:"selfLink"`
}

// Link defines the link attributes such as the link type and target.
type Link struct {
	Type string `json:"type"`
	Href string `json:"href"`
}
