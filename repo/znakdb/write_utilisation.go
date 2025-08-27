package znakdb

import (
	"agregat/domain"
	"agregat/reductor"
	"fmt"
	"time"

	"github.com/upper/db/v4"
)

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
			cis := nextRecords(cises, indexUtilisation, 30000)
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
		IdOrderMarkCodes: int(model.Order),
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
		model.Utilisation = append(model.Utilisation, report.Id)
		for i := range cis {
			km := &domain.UtilisationCodes{
				IdOrderMarkUtilisation: int(report.Id),
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

// получить следующие км
// i номер группы по count штук
// если размер массива меньше count значит последний
// елси размер массива 0 значит больше нет
func nextRecords(arr []*domain.Record, i int, count int) (out []*domain.Record) {
	lenCis := len(arr)
	out = make([]*domain.Record, 0)
	first := i * count
	for i := 0; i < count; i++ {
		index := i + first
		if (index + 1) > lenCis {
			return out
		}
		out = append(out, arr[index])
	}
	return out
}
