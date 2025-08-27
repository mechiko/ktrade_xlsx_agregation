package checkdbg

import (
	"agregat/domain"
	"fmt"

	"go.uber.org/zap"
)

const modError = "pkg:checkdbg"

type Checks struct {
	loger *zap.SugaredLogger
	repo  domain.Repo
}

func NewChecks(loger *zap.SugaredLogger, repo domain.Repo) (*Checks, error) {
	// инициализируем REPO
	// TODO изменить получение путей из конфига
	return &Checks{
		loger: loger,
		repo:  repo,
	}, nil
}

func (c *Checks) Run() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s Run panic %v", modError, r)
		}
	}()
	if err := c.TestDbConfigContact(); err != nil {
		return err
	}
	if err := c.TestDbConfigReleaseMethod(); err != nil {
		return err
	}
	return nil
}
