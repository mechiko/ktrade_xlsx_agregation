package znakdb

import (
	"agregat/domain"
	"errors"
	"fmt"

	"github.com/upper/db/v4"
)

func (z *DbZnak) FindOrders(in []*domain.Record) (err error) {
	sess := z.dbSession
	for i, rec := range in {
		if rec == nil || rec.Cis == nil {
			return fmt.Errorf("record[%d]: CIS is nil", i)
		}
		order := &domain.OrderMarkCodesSerialNumbers{}
		res := sess.Collection("order_mark_codes_serial_numbers").Find("code", rec.Cis.Code).And("status = ?", "Напечатан")
		if err := res.One(order); err != nil {
			if errors.Is(err, db.ErrNoMoreRows) {
				return fmt.Errorf("%s: km %s with state Напечатан not found", modError, rec.Cis.Code)
			}
			return fmt.Errorf("ошибка поиска КМ [%s] в базе %w", rec.Cis.Code, err)
		}
		rec.Order = order.IdOrderMarkCodes
		rec.Serial = order.SerialNumber
	}
	return nil
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
