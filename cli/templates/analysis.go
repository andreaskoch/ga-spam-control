package templates

// PrettyAnalysis contains the detailed template for displaying
// the spam analysis results for a given analytics account.
var PrettyAnalysis = `{{printf "%45s" "Domainname"}} Probability
{{range .SpamDomains -}}
{{printf "%45s" .DomainName}} {{printf "%.2f" .SpamProbability}}
{{end}}`

// QuietAnalysis contains the minimal  template for displaying
// the spam analysis results for a given analytics account.
var QuietAnalysis = `{{range .SpamDomains -}}
{{.DomainName}}
{{end}}`
