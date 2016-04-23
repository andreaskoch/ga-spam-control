package templates

var PrettyStatus = `Global status: {{.OverallStatus}}

Account status:
{{range .Accounts -}}
{{printf "%s (%s)" .AccountName .AccountID | printf "%30s"}}: {{.Status}}
{{end}}`

var QuietStatus = `{{range .Accounts -}}
{{printf "%-10s" .AccountID}} {{printf "%s" .Status}}
{{end}}`
