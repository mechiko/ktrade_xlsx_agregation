package znakdb

import (
	"agregat/domain"
	"errors"
	"fmt"

	"github.com/upper/db/v4"
)

func (z *DbZnak) FindOrders(in []*domain.Record) (errArray []error) {
	errArray = make([]error, 0)
	sess := z.dbSession
	for i, rec := range in {
		if rec == nil || rec.Cis == nil {
			errArray = append(errArray, fmt.Errorf("record[%d]: CIS is nil", i))
			continue
		}
		order := &domain.OrderMarkCodesSerialNumbers{}
		res := sess.Collection("order_mark_codes_serial_numbers").Find("code", rec.Cis.Code)
		// убрал условие чтобы увидеть реальный статус марки .And("status = ?", "Напечатан")
		if err := res.One(order); err != nil {
			if errors.Is(err, db.ErrNoMoreRows) {
				errArray = append(errArray, fmt.Errorf("поиска КМ: [%s] в базе не дал результата %w", rec.Cis.Code))
				continue
			}
			errArray = append(errArray, fmt.Errorf("ошибка поиска КМ [%s] в базе %w", rec.Cis.Code, err))
			continue
		}
		if order.Status != "Напечатан" {
			errArray = append(errArray, fmt.Errorf("поиск КМ: [%s - %s] нужен - Напечатан", rec.Cis.Code, order.Status))
		}
		rec.Order = order.IdOrderMarkCodes
		rec.Serial = order.SerialNumber
	}
	return errArray
}

// func (z *DbZnak) Order(id int64) (sl []*dbznak.OrderMarkCodes, err error) {
// 	sess := z.dbSession
// 	defer sess.Close()

// 	err = sess.Collection("order_mark_codes").Find("cis_type like ? and archive <> 1", cisType).OrderBy("id desc").All(&sl)
// 	if err != nil {
// 		return nil, fmt.Errorf("db:znak upper CisTypeOrders %w", err)
// 	}
// 	return sl, nil
// }
