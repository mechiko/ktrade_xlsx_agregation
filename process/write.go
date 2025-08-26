package process

import (
	"agregat/reductor"
	"agregat/repo/znakdb"
	"fmt"

	"github.com/mechiko/dbscan"
)

func (p *process) WriteUtilisation() (err error) {
	info := p.repo.Info(dbscan.TrueZnak)
	if info == nil {
		return fmt.Errorf("базы 4z не найдено")
	}
	db, err := znakdb.New(info, dbscan.TrueZnak)
	if err != nil {
		return fmt.Errorf("error open 4z db")
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			if err != nil {
				// keep original op error and append close error
				err = fmt.Errorf("%w; close error: %v", err, cerr)
			} else {
				err = cerr
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
	return
}
