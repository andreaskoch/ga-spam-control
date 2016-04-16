package apimodel

type Filter struct {
	AccountID      string
	ID             string
	Kind           string
	Name           string
	Type           string
	Link           string
	ExcludeDetails FilterDetail
}

type FilterDetail struct {
	Kind            string
	Field           string
	MatchType       string
	ExpressionValue string
	CaseSensitive   bool
}
