package process

import (
	"agregat/htmltmpl"
	"fmt"
	"os"
	"path/filepath"
)

// возвращает имя файла отчета
func (p *process) SuccessHtml(outDir string) (string, error) {
	htmString, err := htmltmpl.NewTemplate().SuccessHTML(p)
	if err != nil {
		return "", fmt.Errorf("ошибка template %w", err)
	}
	fileHtml := "report_" + p.NameFileWithoutExt
	fileHtml = filepath.Join(outDir, fileHtml) + ".html"
	file, err := os.Create(fileHtml)
	if err != nil {
		return "", fmt.Errorf("ошибка создания файла %w", err)
	}
	defer file.Close() // Ensure the file is closed when the function exits

	// Write the string content to the file
	_, err = file.WriteString(string(htmString))
	if err != nil {
		return "", fmt.Errorf("ошибка записи %w", err)
	}
	return fileHtml, nil
}

// возвращает имя файла отчета
func (p *process) ErrorHtml(outDir string) (string, error) {
	htmString, err := htmltmpl.NewTemplate().ErrorHTML(p)
	if err != nil {
		return "", fmt.Errorf("ошибка template %w", err)
	}
	fileHtml := "error_" + p.NameFileWithoutExt
	fileHtml = filepath.Join(outDir, fileHtml) + ".html"
	file, err := os.Create(fileHtml)
	if err != nil {
		return "", fmt.Errorf("ошибка создания файла %w", err)
	}
	defer file.Close() // Ensure the file is closed when the function exits

	// Write the string content to the file
	_, err = file.WriteString(string(htmString))
	if err != nil {
		return "", fmt.Errorf("ошибка записи %w", err)
	}
	return fileHtml, nil
}
