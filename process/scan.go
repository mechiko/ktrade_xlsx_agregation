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
		gtin := strings.TrimSpace(rec.Cis.Gtin)
		code := strings.TrimSpace(rec.Cis.Code)
		korob := strings.TrimSpace(rec.Korob)
		palet := strings.TrimSpace(rec.Palet)
		produced := rec.Produced.Format("02.01.2006")
		if _, ok := p.RecordsMap[cis]; !ok {
			p.RecordsMap[cis] = rec
		} else {
			return fmt.Errorf("double KM cis %s %d", cis, iRec)
		}
		p.KM[cis] = rec.Cis
		p.arrKM = append(p.arrKM, code)
		if _, ok := p.Koroba[korob]; !ok {
			p.Koroba[rec.Korob] = &domain.Korob{
				KITU:     korob,
				GTIN:     gtin,
				Km:       make([]string, 0),
				Produced: produced,
			}
			p.KorobaKeys = append(p.KorobaKeys, korob)
		}
		if p.Koroba[korob].GTIN != gtin {
			return fmt.Errorf("gtin cis %s not equal korob %s", p.Koroba[korob].GTIN, gtin)
		}
		if p.Koroba[korob].Produced != produced {
			return fmt.Errorf("produced cis %s not equal korob %s", p.Koroba[korob].Produced, produced)
		}
		p.Koroba[rec.Korob].Km = append(p.Koroba[rec.Korob].Km, cis)
		if _, ok := p.Palet[rec.Palet]; !ok {
			p.Palet[palet] = &domain.Palet{
				KITU:     palet,
				GTIN:     gtin,
				Korobs:   make([]string, 0),
				Produced: produced,
			}
			if _, ok := p.PaletByDateProduce[produced]; !ok {
				p.PaletByDateProduce[produced] = make([]string, 0)
			}
			p.PaletByDateProduce[produced] = append(p.PaletByDateProduce[produced], palet)
		}
		if p.Palet[palet].GTIN != gtin {
			return fmt.Errorf("gtin korob %s not equal palet %s", p.Koroba[korob].GTIN, gtin)
		}
		if p.Palet[palet].Produced != produced {
			return fmt.Errorf("produced korob %s not equal palet %s", p.Koroba[korob].Produced, produced)
		}
		p.Palet[palet].Korobs = append(p.Palet[palet].Korobs, korob)
	}
	p.ListKoroba = make([][]string, 0)
	keysKorob := make([]string, 0, len(p.Koroba))
	for k := range p.Koroba {
		keysKorob = append(keysKorob, k)
	}
	slices.Sort(keysKorob)
	for _, key := range keysKorob {
		for _, cis := range p.Koroba[key].Km {
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
		keys := make([]string, 0, len(p.Palet[key].Korobs))
		for _, k := range p.Palet[key].Korobs {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, kk := range keys {
			r := []string{key, kk}
			p.ListPalet = append(p.ListPalet, r)
		}
	}
	return nil
}
