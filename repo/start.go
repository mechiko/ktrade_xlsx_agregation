package repo

import (
	"context"
	"fmt"
)

func (r *Repository) Run(ctx context.Context) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
		}
	}()
	// ожидаем завершения контекста
	<-ctx.Done()
	r.logger.Infof("завершаем работу репозитория по контексту")
	return nil
}
