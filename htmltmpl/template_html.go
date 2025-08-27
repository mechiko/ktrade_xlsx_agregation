package htmltmpl

import (
	"bytes"
	"fmt"
	"html/template"
)

func (tt templateString) tmplMustText(tmpl string, tmplName string, data interface{}, f template.FuncMap) (ss []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			ss = nil
			err = fmt.Errorf("panic templateString %v", r)
		}
	}()

	var buf bytes.Buffer
	fncMap := funcMapText
	if f != nil {
		fncMap = f
	}
	t := template.Must(template.New(tmplName).Funcs(fncMap).Parse(tmpl))
	err = t.ExecuteTemplate(&buf, tmplName, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
