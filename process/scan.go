package process

import (
	"agregat/domain"
	"agregat/repo/znakdb"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/mechiko/dbscan"
	"github.com/upper/db/v4"
)

func (p *process) ScanRecords() (err error) {
	p.KMErrors = make([]string, 0)
	info := p.repo.Info(dbscan.TrueZnak)
	if info == nil {
		return fmt.Errorf("базы 4z не найдено")
	}
	dbZnak, err := znakdb.New(info, dbscan.TrueZnak)
	if err != nil {
		return fmt.Errorf("open znak db: %w", err)
	}
	defer func() {
		if cerr := dbZnak.Close(); cerr != nil {
			if err != nil {
				// keep original op error and append close error
				err = fmt.Errorf("%w; close error: %v", err, cerr)
			} else {
				err = cerr
			}
		}
	}()
	p.Guide, err = dbZnak.Guide()
	if err != nil {
		return fmt.Errorf("error get guide %w", err)
	}
	if err := dbZnak.FindOrders(p.Records); len(err) != 0 {
		for _, v := range err {
			p.KMErrors = append(p.KMErrors, v.Error())
		}
		return fmt.Errorf("error scan km contains errors %d", len(err))
	}
	for _, rec := range p.Records {
		pal := strings.TrimSpace(rec.Palet)
		plt, err := dbZnak.FindPallet(pal)
		if err != nil && !errors.Is(err, db.ErrNoMoreRows) {
			return fmt.Errorf("find palet %s: %w", pal, err)
		}
		if err == nil && plt != nil {
			return fmt.Errorf("palet is present %s created %s id %v", plt["unit_serial_number"], plt["create_date"], plt["id"])
		}
		kor := strings.TrimSpace(rec.Korob)
		krb, err := dbZnak.FindPallet(kor)
		if err != nil && !errors.Is(err, db.ErrNoMoreRows) {
			return fmt.Errorf("find korob %s: %w", kor, err)
		}
		if err == nil && krb != nil {
			return fmt.Errorf("korob is present %s created %s id %v", krb["unit_serial_number"], krb["create_date"], krb["id"])
		}
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
		if _, ok := p.Palet[palet]; !ok {
			p.Palet[palet] = &domain.Palet{
				KITU:     palet,
				GTIN:     gtin,
				Korobs:   make([]string, 0),
				Produced: produced,
			}
			orderKey := fmt.Sprintf("%s:%s", gtin, produced)
			if _, ok := p.PaletByDateProduce[orderKey]; !ok {
				p.PaletByDateProduce[orderKey] = make([]string, 0)
			}
			p.PaletByDateProduce[orderKey] = append(p.PaletByDateProduce[orderKey], palet)
		}
		if p.Palet[palet].GTIN != gtin {
			return fmt.Errorf("palet %s: GTIN mismatch (have %s, got %s)", palet, p.Palet[palet].GTIN, gtin)
		}
		if p.Palet[palet].Produced != produced {
			return fmt.Errorf("palet %s: produced mismatch (have %s, got %s)", palet, p.Palet[palet].Produced, produced)
		}
		p.arrKM = append(p.arrKM, code)
		if _, ok := p.Koroba[korob]; !ok {
			p.Koroba[korob] = &domain.Korob{
				KITU:     korob,
				GTIN:     gtin,
				Km:       make([]string, 0),
				Produced: produced,
			}
			p.KorobaKeys = append(p.KorobaKeys, korob)
			p.Palet[palet].Korobs = append(p.Palet[palet].Korobs, korob)
		}
		if p.Koroba[korob].GTIN != gtin {
			return fmt.Errorf("korob %s: GTIN mismatch (have %s, got %s)", korob, p.Koroba[korob].GTIN, gtin)
		}
		if p.Koroba[korob].Produced != produced {
			return fmt.Errorf("korob %s: produced mismatch (have %s, got %s)", korob, p.Koroba[korob].Produced, produced)
		}
		p.Koroba[korob].Km = append(p.Koroba[korob].Km, cis)
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
