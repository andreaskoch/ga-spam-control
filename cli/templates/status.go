package templates

var Status = `{{range .Accounts -}}
{{printf "%s (%s)" .AccountName .AccountID | printf "%30s"}}: {{.Status}}
{{end}}`
