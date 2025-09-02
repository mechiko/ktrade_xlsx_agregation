package htmltmpl

import (
	"bytes"
	"fmt"
	"html/template"
)

func (tt templateString) tmplHtml(tmpl string, tmplName string, data interface{}, f template.FuncMap) (ss []byte, err error) {
	var buf bytes.Buffer
	fncMap := funcMapHtml
	if f != nil {
		fncMap = f
	}
	t, err := template.New(tmplName).Funcs(fncMap).Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("tmplHtml parse error %w", err)
	}
	err = t.ExecuteTemplate(&buf, tmplName, data)
	if err != nil {
		return nil, fmt.Errorf("tmplHtml execute error %w", err)
	}
	return buf.Bytes(), nil
}
