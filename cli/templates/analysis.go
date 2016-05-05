package templates

// PrettyAnalysis contains the detailed template for displaying
// the spam analysis results for a given analytics account.
var PrettyAnalysis = `{{if .SpamDomains -}}
{{printf "%45s" "Domainname"}} Probability
{{range .SpamDomains -}}
{{printf "%45s" .DomainName}} {{printf "%.2f" .SpamProbability}}
{{end -}}
{{else -}}
No referrer spam domains detected
{{- end}}`

// QuietAnalysis contains the minimal  template for displaying
// the spam analysis results for a given analytics account.
var QuietAnalysis = `{{if .SpamDomains -}}
{{range .SpamDomains -}}
{{.DomainName}}
{{end}}
{{- end}}`
