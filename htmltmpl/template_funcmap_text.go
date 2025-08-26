package htmltmpl

import (
	"bytes"
	"encoding/xml"
	"text/template"
	"time"
)

var funcMapText = template.FuncMap{
	// The name "inc" is what the function will be called in the template text.
	"dt": func(i time.Time) string {
		return i.Format("02.01.2006")
	},
	"escape": func(s string) string {
		var sh bytes.Buffer
		xml.Escape(&sh, []byte(s))
		return sh.String()
	},
}
