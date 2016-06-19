package templates

// PrettyAnalysis contains the detailed template for displaying
// the spam analysis results for a given analytics account.
var PrettyAnalysis = `{{if .Domains -}}
{{range .Domains -}}
{{.DomainName}}
{{end -}}
{{else -}}
No new referrer spam domains detected
{{end}}`

// QuietAnalysis contains the minimal  template for displaying
// the spam analysis results for a given analytics account.
var QuietAnalysis = `{{if .Domains -}}
{{range .Domains -}}
{{.DomainName}}
{{end -}}
{{end}}`
