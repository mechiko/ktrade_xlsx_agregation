package repo

import (
	"errors"
	"fmt"

	"github.com/mechiko/dbscan"
)

// при инициализации приложения этот метод вызывается однажды и прописывает объект доступа
// к базе данных, далее проверяет версию БД возможна ошибка и нужно выходить из приложения
func (r *Repository) prepareSelf() (err error) {
	defer func() {
		if rr := recover(); rr != nil {
			err = fmt.Errorf("%s panic %v", modError, rr)
		}
	}()
	// lock создает объект базы при этом она открывается и мы ответственны о ее закрытии и unlock
	self, err := r.Lock(dbscan.Other)
	if err != nil {
		return fmt.Errorf("%s lock %v", modError, err)
	}
	// only Close/Unlock if Lock held the mutex (i.e., returned a non-nil DB)
	if self != nil {
		defer func() {
			var dErr error
			if errClose := self.Close(); errClose != nil {
				r.logger.Errorf("%s self.Close %v", modError, errClose)
				dErr = errors.Join(dErr, fmt.Errorf("%s selfdb close: %w", modError, errClose))
			}
			if errUnLock := r.Unlock(dbscan.Other); errUnLock != nil {
				r.logger.Errorf("%s unlock error %v", modError, errUnLock)
				dErr = errors.Join(dErr, fmt.Errorf("%s selfdb unlock: %w", modError, errUnLock))
			}
			if dErr != nil {
				err = errors.Join(err, dErr)
			}
		}()
	}

	info := r.Info(dbscan.Other)
	if info == nil {
		return fmt.Errorf("%s get info: nil", modError)
	}
	if !info.Exists {
		return fmt.Errorf("%s get self db error!", modError)
	}
	return nil
}
