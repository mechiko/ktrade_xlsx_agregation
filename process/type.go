package process

import (
	"agregat/domain"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mechiko/utility"
)

type process struct {
	repo               domain.Repo
	NameFileWithoutExt string
	File               string
	Records            []*domain.Record
	RecordsMap         map[string]*domain.Record
	Koroba             map[string]*domain.Korob
	Palet              map[string]*domain.Palet
	PaletByDateProduce map[string][]string
	KM                 map[string]*utility.CisInfo
	KMErrors           []string
	arrKM              []string
	ListKoroba         [][]string
	ListPalet          [][]string
	KorobaKeys         []string
	SerialNumberType   string
	CisType            string
	ContactPerson      string
	ReleaseMethodType  string
	CreateMethodType   string
	PaymentType        string
	TemplateId         string
	Utilisation        map[string]*UtilisationReport
	Guide              map[string]*domain.ProductGuides
}

type UtilisationReport struct {
	Order    int64
	Prod     time.Time
	Exp      time.Time
	ReportID int64
	KM       []*domain.Record
}

func (u *UtilisationReport) String() string {
	out := fmt.Sprintf("%d:%s", u.Order, u.Prod.Local().Format("2006.01.02"))
	return out
}

func New(file string, repo domain.Repo) (*process, error) {
	if repo == nil {
		return nil, fmt.Errorf("repo is nil")
	}
	if file == "" {
		return nil, fmt.Errorf("file name empty")
	}
	if !utility.PathOrFileExists(file) {
		return nil, fmt.Errorf("file not found")
	}
	name := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	p := &process{
		repo:               repo,
		File:               file,
		NameFileWithoutExt: name,
		Koroba:             make(map[string]*domain.Korob),
		Palet:              make(map[string]*domain.Palet),
		KM:                 make(map[string]*utility.CisInfo),
		arrKM:              make([]string, 0),
		Records:            make([]*domain.Record, 0),
		RecordsMap:         make(map[string]*domain.Record),
		Utilisation:        make(map[string]*UtilisationReport),
		PaletByDateProduce: make(map[string][]string),
	}
	return p, nil
}
