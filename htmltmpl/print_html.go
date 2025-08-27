package htmltmpl

import (
	_ "embed"
	"fmt"
)

//go:embed tmplHtml.html
var tmplHtml string

func (tt *templateString) StringHTML(model interface{}) (bts []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic templateString %v", r)
		}
	}()

	tmplName := "html"
	// вызов шаблона в него передаем имя шаблона как имя файла шаблона
	if result, err := tt.tmplMustText(tmplHtml, tmplName, model, nil); err != nil {
		return bts, err
	} else {
		return result, err
	}
}
