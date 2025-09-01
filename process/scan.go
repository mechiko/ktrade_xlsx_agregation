package process

import (
	"agregat/domain"
	"agregat/repo/znakdb"
	"fmt"
	"slices"
	"strings"

	"github.com/mechiko/dbscan"
)

func (p *process) ScanRecords() (err error) {
	p.KMErrors = make([]string, 0)
	info := p.repo.Info(dbscan.TrueZnak)
	if info == nil {
		return fmt.Errorf("базы 4z не найдено")
	}
	db, err := znakdb.New(info, dbscan.TrueZnak)
	if err != nil {
		return fmt.Errorf("open znak db: %w", err)
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
	if err := db.FindOrders(p.Records); len(err) != 0 {
		for _, v := range err {
			p.KMErrors = append(p.KMErrors, v.Error())
		}
		return fmt.Errorf("error scan km contains errors %d", len(err))
	}
	for _, rec := range p.Records {
		ur := &UtilisationReport{
			Order: rec.Order,
			Prod:  rec.Produced,
			Exp:   rec.Expired,
			KM:    make([]*domain.Record, 0),
		}
		urStr := ur.String()
		if _, ok := p.Utilisation[urStr]; !ok {
			p.Utilisation[urStr] = ur
		}
		p.Utilisation[urStr].KM = append(p.Utilisation[urStr].KM, rec)
	}
	return
}

func (p *process) ScanPalet() (err error) {
	for iRec, rec := range p.Records {
		cis := strings.TrimSpace(rec.Cis.Cis)
		if _, ok := p.RecordsMap[cis]; !ok {
			p.RecordsMap[cis] = rec
		} else {
			return fmt.Errorf("double KM cis %s %d", cis, iRec)
		}
		p.KM[cis] = rec.Cis
		p.arrKM = append(p.arrKM, rec.Cis.Code)
		if _, ok := p.Koroba[rec.Korob]; !ok {
			p.Koroba[rec.Korob] = make([]string, 0)
			p.KorobaKeys = append(p.KorobaKeys, rec.Korob)
		}
		p.Koroba[rec.Korob] = append(p.Koroba[rec.Korob], cis)
		if _, ok := p.Palet[rec.Palet]; !ok {
			p.Palet[rec.Palet] = make(map[string]string)
		}
		p.Palet[rec.Palet][rec.Korob] = rec.Korob
	}
	p.ListKoroba = make([][]string, 0)
	keysKorob := make([]string, 0, len(p.Koroba))
	for k := range p.Koroba {
		keysKorob = append(keysKorob, k)
	}
	slices.Sort(keysKorob)
	for _, key := range keysKorob {
		for _, cis := range p.Koroba[key] {
			r := []string{key, cis}
			p.ListKoroba = append(p.ListKoroba, r)
		}
	}
	p.ListPalet = make([][]string, 0)
	keysPalet := make([]string, 0, len(p.Palet))
	for k := range p.Palet {
		keysPalet = append(keysPalet, k)
	}
	slices.Sort(keysPalet)
	for _, key := range keysPalet {
		keys := make([]string, 0, len(p.Palet[key]))
		for k := range p.Palet[key] {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, kk := range keys {
			r := []string{key, kk}
			p.ListPalet = append(p.ListPalet, r)
		}
	}

	for _, key := range keysPalet {
		for kk := range p.Palet[key] {
			if kk != "" {
				date, err := p.findKorob(kk)
				if err != nil {
					return fmt.Errorf("scanpalet order palet by date %w", err)
				}
				if _, ok := p.PaletByDateProduce[date]; !ok {
					p.PaletByDateProduce[date] = make([]string, 0)
				}
				p.PaletByDateProduce[date] = append(p.PaletByDateProduce[date], key)
			} else {
				return fmt.Errorf("scanpalet order palet by date empty korob in palet %s", key)
			}
			break
		}
	}
	return nil
}

func (p *process) findKorob(korob string) (date string, err error) {
	recs, ok := p.Koroba[korob]
	if !ok {
		return "", fmt.Errorf("find korob date нет такого короба %s", korob)
	}
	if len(recs) > 0 {
		cis := strings.TrimSpace(recs[0])
		rec, ok := p.RecordsMap[cis]
		if !ok {
			return "", fmt.Errorf("find korob date нет такой KM %s", cis)
		}
		date = rec.Produced.Format("02.01.2006")
		return date, nil
	}
	return "", fmt.Errorf("find korob date нет марок в коробе %s", korob)
}
