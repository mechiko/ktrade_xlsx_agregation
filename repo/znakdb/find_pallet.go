package znakdb

import (
	"fmt"
)

// возвращает map[string]interface{}
// возвращает ошибку когда не найдено и когда реально ошибка
// если что, такая ошибка может быть errors.Is(err, db.ErrNoMoreRows)
func (z *DbZnak) FindPallet(number string) (pallet map[string]interface{}, err error) {
	sess := z.dbSession
	defer func() {
		if err != nil {
			if errClose := sess.Close(); errClose != nil {
				err = fmt.Errorf("%w%w", errClose, err)
			}
		} else {
			err = sess.Close()
		}
	}()

	pallet = make(map[string]interface{})
	res := sess.Collection("order_mark_aggregation").Find("unit_serial_number", number)
	if err := res.One(&pallet); err != nil {
		// if errors.Is(err, db.ErrNoMoreRows) {
		// }
		return pallet, err
	}
	return pallet, nil
}
