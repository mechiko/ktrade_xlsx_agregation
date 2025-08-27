package process

import (
	"agregat/domain"
	"fmt"

	"github.com/xuri/excelize/v2"
)

func (k *process) ReadXlsx(file string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic ReadXlsx %v", r)
		}
	}()

	f, err := excelize.OpenFile(file)
	if err != nil {
		return fmt.Errorf("open xlsx error %w", err)
	}
	defer func() {
		// Close the spreadsheet.
		if errr := f.Close(); errr != nil {
			err = errr
		}
	}()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return fmt.Errorf("len sheet wrong")
	}
	currentSheet := sheets[0]

	// Получить все строки в Sheet1
	rows, err := f.GetRows(currentSheet)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for rowNumber, row := range rows {
		if domain.IsRecord(row) {
			// берем запись
			rec, err := domain.NewRecord(row)
			if err != nil {
				return fmt.Errorf("error read xlsx row %d %v", rowNumber+1, err)
			}
			k.Records = append(k.Records, rec)
		}
	}
	return nil
}
