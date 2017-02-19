package message

import (
	"bytes"
	"html/template"

	"github.com/shihanng/oss-manager/db"
)

const updateMessage = `
{{- if len .Versions | lt 1 -}}
  Updates for {{.Name}} are available:
{{else -}}
  An update for {{.Name}} is available:
{{end}}
{{range $index, $element := .Versions}}
  {{- printf "  - %s\n" $element}}
{{- end}}
 {{.URL}}
`

const listMessage = `
{{- range $i, $project := .}}
{{- printf "%s:\n" $project.Name}}
  {{- range $i, $version := $project.Versions}}
  {{- printf " %s" $version}}
{{end}}
  {{- printf " %s\n\n" $project.URL}}
{{- end -}}
`

var (
	tmplUpdate *template.Template
	tmplList   *template.Template
)

func init() {
	tmplUpdate = template.Must(template.New("update").Parse(updateMessage))
	tmplList = template.Must(template.New("list").Parse(listMessage))
}

func ForUpdate(p db.Project) (string, error) {
	var out bytes.Buffer

	err := tmplUpdate.Execute(&out, p)
	if err != nil {
		return "", nil
	}

	return out.String(), nil
}

func ForList(p []db.Project) (string, error) {
	var out bytes.Buffer

	err := tmplList.Execute(&out, p)
	if err != nil {
		return "", nil
	}

	return out.String(), nil
}
