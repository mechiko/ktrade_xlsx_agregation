package repo

import (
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
	defer func() {
		if errClose := self.Close(); errClose != nil {
			r.logger.Errorf("%s self.Close %v", modError, errClose)
			err = fmt.Errorf("%s error selfdb close %w", modError, errClose)
		}
		if errUnLock := r.Unlock(dbscan.Other); errUnLock != nil {
			r.logger.Errorf("%s unlock error %v", modError, errUnLock)
			err = fmt.Errorf("%s error selfdb close", modError)
		}
	}()

	info := r.Info(dbscan.Other)
	if info == nil {
		return fmt.Errorf("%s get info %v", modError, err)
	}
	if !info.Exists {
		return fmt.Errorf("%s get self db error!", modError)
	}
	return nil
}
