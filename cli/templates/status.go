package templates

// PrettyStatus contains the detailed template
// for displaying the status of multiple analytics
// accounts.
var PrettyStatus = `Global status: {{.OverallStatus}}

Account status:
{{range .Accounts -}}
{{printf "%s (%s)" .AccountName .AccountID | printf "%30s"}}: {{.Status}}
{{end}}`

// QuietStatus contains a minimal template for
// displaying the status of multiple analytics accounts.
var QuietStatus = `{{range .Accounts -}}
{{printf "%-10s" .AccountID}} {{printf "%s" .Status}}
{{end}}`
