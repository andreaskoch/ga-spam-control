package templates

// PrettyUpdate contains the detailed template
// for displaying domain list update results.
var PrettyUpdate = `Total: {{.Statistics.Total}} | Added (+): {{.Statistics.Added}} | Removed (-): {{.Statistics.Removed}}  | Unchanged: {{.Statistics.Unchanged}}

{{range .Domains -}}
{{if not .UpdateType.IsUnchanged -}}
{{printf "  %s  %s" .UpdateType .Domainname}}
{{end -}}
{{end}}`

// QuietUpdate contains the minimal template
// for displaying domain list update results.
var QuietUpdate = `{{range .Domains -}}
{{printf "  %s  %s" .UpdateType .Domainname}}
{{end}}`
