// Package templates contains all text-ui templates for the command-line
// interface.
package templates

// PrettyStatus contains the detailed template
// for displaying the status of multiple analytics
// accounts.
var PrettyStatus = `Known spam domains: {{.KnownSpamDomains}}

{{range .Accounts -}}
{{printf "%s (%s)" .AccountName .AccountID | printf "%30s"}}: {{printf "%5s" .Status}} {{printf "%d/%d" .Status.UpToDateFilters .Status.TotalFilters | printf "%8s"}}
{{end}}`

// QuietStatus contains a minimal template for
// displaying the status of multiple analytics accounts.
var QuietStatus = `{{range .Accounts -}}
{{printf "%-10s" .AccountID}} {{.Status}}
{{end}}`
