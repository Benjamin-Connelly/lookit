package templates

import (
	"embed"
	"html/template"
	"strings"
)

//go:embed *.html
var templateFS embed.FS

var Templates *template.Template

func init() {
	funcMap := template.FuncMap{
		"add":        func(a, b int) int { return a + b },
		"repeat":     strings.Repeat,
		"trimPrefix": strings.TrimPrefix,
	}
	Templates = template.Must(template.New("").Funcs(funcMap).ParseFS(templateFS, "*.html"))
}
