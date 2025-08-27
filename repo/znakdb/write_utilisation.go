package znakdb

import (
	"agregat/domain"
	"agregat/reductor"
	"fmt"
	"time"

	"github.com/upper/db/v4"
)

const maxCountPerReportUtilisation = 300000

// запись отчета нанесения
func (z *DbZnak) WriteUtilisation(cises []*domain.Record, model *reductor.Model, prod, exp time.Time) (rid int64, err error) {
	defer func() {
		if errRecover := recover(); errRecover != nil {
			if err != nil {
				err = fmt.Errorf("%s %v %w", modError, errRecover, err)
			} else {
				err = fmt.Errorf("%s %v", modError, errRecover)
			}
		}
	}()

	sess := z.dbSession
	err = sess.Tx(func(tx db.Session) error {
		indexUtilisation := 0
		for {
			cis := nextRecords(cises, indexUtilisation)
			indexUtilisation++
			if len(cis) == 0 {
				// больше нет км
				break
			}
			if ri, err := z.writeUtilisation(tx, cis, model, prod, exp); err != nil {
				return err
			} else {
				rid = ri
			}
		}
		return nil
	})
	return rid, err
}

func (z *DbZnak) writeUtilisation(tx db.Session, cis []*domain.Record, model *reductor.Model, prod, exp time.Time) (rid int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic %v", r)
		}
	}()
	report := &domain.Utilisation{
		CreateDate:       time.Now().Local().Format("2006.01.02 15:04:05"),
		PrimaryDocDate:   time.Now().Local().Format("2006.01.02 15:04:05"),
		IdOrderMarkCodes: model.Order,
		ProductionDate:   prod.Local().Format("2006.01.02"),
		ExpirationDate:   exp.Local().Format("2006.01.02"),
		UsageType:        "Нанесение КМ подтверждено",
		Inn:              model.Inn,
		Kpp:              model.Kpp,
		Version:          "1",
		State:            "Создан",
		Status:           "Не проведён",
		Quantity:         fmt.Sprintf("%d", len(cis)),
		ReportId:         "",
		Archive:          0,
		AlcVolume:        "",
	}
	if err := tx.Collection("order_mark_utilisation").InsertReturning(report); err != nil {
		return 0, err
	} else {
		// model.Utilisation = append(model.Utilisation, report.Id)
		for i := range cis {
			if cis[i] == nil || cis[i].Cis == nil {
				return 0, fmt.Errorf("cis[%d]: nil record or CIS", i)
			}
			km := &domain.UtilisationCodes{
				IdOrderMarkUtilisation: report.Id,
				SerialNumber:           cis[i].Serial,
				Code:                   cis[i].Cis.Code,
				Status:                 "Нанесён",
			}
			if _, err := tx.Collection("order_mark_utilisation_codes").Insert(km); err != nil {
				return 0, err
			}
		}
	}
	return report.Id, nil
}

// nextRecords returns a batch of records starting from startIndex
// Returns empty slice when no more records are available
func nextRecords(arr []*domain.Record, index int) []*domain.Record {
	startIndex := index * maxCountPerReportUtilisation
	if startIndex >= len(arr) {
		return []*domain.Record{}
	}

	endIndex := startIndex + maxCountPerReportUtilisation
	if endIndex > len(arr) {
		endIndex = len(arr)
	}

	return arr[startIndex:endIndex]
}
