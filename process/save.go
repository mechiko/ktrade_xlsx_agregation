package process

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (p *process) Save(out string) error {
	fileNameKoroba := "koroba_" + p.NameFileWithoutExt + ".csv"
	fileNameKoroba = filepath.Join(out, fileNameKoroba)
	if err := saveTxt(fileNameKoroba, p.ListKoroba); err != nil {
		return fmt.Errorf("error write file koroba %w", err)
	}
	fileNamePalet := "palet_" + p.NameFileWithoutExt + ".csv"
	fileNamePalet = filepath.Join(out, fileNamePalet)
	if err := saveTxt(fileNamePalet, p.ListPalet); err != nil {
		return fmt.Errorf("error write file palet %w", err)
	}
	return nil
}

func saveTxt(name string, data [][]string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, row := range data {
		_, err = fmt.Fprintf(file, "%s\n", strings.Join(row, "\t"))
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}
	return nil
}

func saveKM(name string, data []string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, row := range data {
		_, err = fmt.Fprintf(file, "%s\n", row)
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}
	return nil
}

func NewWriter(w io.Writer) (writer *csv.Writer) {
	writer = csv.NewWriter(w)
	writer.Comma = '\t'
	return
}

func saveCsvCustom(name string, data [][]string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '\t'
	writer.WriteAll(data) // calls Flush internally
	return writer.Error()
}
