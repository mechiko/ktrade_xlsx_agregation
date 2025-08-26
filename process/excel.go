package process

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func ExcelSheets(file string) (out []string, err error) {
	out = make([]string, 0)
	f, err := excelize.OpenFile(file)
	if err != nil {
		return out, fmt.Errorf("%w", err)
	}
	defer func() {
		// Close the spreadsheet.
		err = f.Close()
	}()
	list := f.GetSheetList()
	for _, sheet := range list {
		rows, err := f.GetRows(sheet)
		if err != nil {
			return out, fmt.Errorf("%w", err)
		}
		for _, row := range rows {
			if len(row) > 0 {
				out = append(out, fmt.Sprintf("%s:%s", sheet, row[0]))
			}
		}
	}
	return out, nil
}

func Excel(file string) (out []string, err error) {
	out = make([]string, 0)
	f, err := excelize.OpenFile(file)
	if err != nil {
		return out, fmt.Errorf("%w", err)
	}
	defer func() {
		// Close the spreadsheet.
		err = f.Close()
	}()
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return out, fmt.Errorf("%w", err)
	}
	for _, row := range rows {
		if len(row) > 0 {
			out = append(out, row[0])
		}
	}
	return out, nil
}
