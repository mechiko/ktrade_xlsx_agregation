package htmltmpl

import (
	"bytes"
	"encoding/xml"
	"strings"
	"text/template"
	"time"
)

var funcMapHtml = template.FuncMap{
	// The name "inc" is what the function will be called in the template text.
	"dt": func(i time.Time) string {
		return i.Format("02.01.2006")
	},
	"escape": func(s string) string {
		var sh bytes.Buffer
		xml.Escape(&sh, []byte(s))
		return sh.String()
	},
	"gtin": func(s string) string {
		ar := strings.Split(s, ":")
		if len(ar) > 0 {
			return ar[0]
		}
		return ""
	},
	"produced": func(s string) string {
		ar := strings.Split(s, ":")
		if len(ar) > 1 {
			return ar[1]
		}
		return ""
	},
}
