package znakdb

import (
	"fmt"
)

// вызывается каждый раз при создании объекта
func (r *DbZnak) Check() (err error) {
	defer func() {
		if rr := recover(); rr != nil {
			err = fmt.Errorf("%s check panic %v", modError, rr)
		}
	}()

	if !r.dbInfo.Exists {
		return fmt.Errorf("dbConfig dbInfo.Exists false")
	}
	r.dbSession, err = r.dbInfo.Connect()
	if err != nil {
		return fmt.Errorf("%s check ошибка подключения к БД %w", modError, err)
	}
	return nil
}
