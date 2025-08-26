package process

import (
	"agregat/reductor"
	"agregat/repo/znakdb"
	"fmt"

	"github.com/mechiko/dbscan"
)

func (p *process) WriteUtilisation() error {
	info := p.repo.Info(dbscan.TrueZnak)
	if info == nil {
		return fmt.Errorf("базы 4z не найдено")
	}
	db, err := znakdb.New(info, dbscan.Config)
	if err != nil {
		return fmt.Errorf("error open config.db")
	}
	defer func() {
		if errclose := db.Close(); errclose != nil {
			if err != nil {
				err = fmt.Errorf("%v %w", errclose, err)
			} else {
				err = fmt.Errorf("%v", errclose)
			}
		}
	}()
	model := reductor.Instance().Model("")
	for _, ur := range p.Utilisation {
		model.Order = ur.Order
		if rid, err := db.WriteUtilisation(ur.KM, &model, ur.Prod, ur.Exp); err != nil {
			return fmt.Errorf("ошибка записи отчета нанесения %w", err)
		} else {
			ur.ReportID = rid
		}
	}
	reductor.Instance().SetModel("", &model)
	return nil
}
