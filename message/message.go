package message

import (
	"bytes"
	"html/template"

	"github.com/shihanng/oss-manager/db"
)

const updateMessage = `{{if len .Versions | lt 1 -}}
Updates for {{.Name}} are available:
{{else -}}
An update for {{.Name}} is available:
{{end}}
{{range $index, $element := .Versions}}{{printf "  - %s\n" $element}}{{end}}
 {{.URL}}
`

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("update").Parse(updateMessage))
}

func ForUpdate(p db.Project) (string, error) {
	var out bytes.Buffer

	err := tmpl.Execute(&out, p)
	if err != nil {
		return "", nil
	}

	return out.String(), nil
}
