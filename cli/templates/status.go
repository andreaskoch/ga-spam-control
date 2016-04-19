package templates

var Status = `Global status: {{.OverallStatus}}

Account status:
{{range .Accounts -}}
{{printf "%s (%s)" .AccountName .AccountID | printf "%30s"}}: {{.Status}}
{{end}}`
