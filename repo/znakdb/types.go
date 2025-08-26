package znakdb

import (
	_ "embed"
	"fmt"

	"github.com/mechiko/dbscan"
	"github.com/upper/db/v4"
	"go.uber.org/zap"
)

const modError = "selfdb"

type DbZnak struct {
	logger    *zap.SugaredLogger
	dbSession db.Session // открытый хэндл тут
	dbInfo    *dbscan.DbInfo
	infoType  dbscan.DbInfoType
	version   int64
}

func New(info *dbscan.DbInfo, infoType dbscan.DbInfoType) (*DbZnak, error) {
	db := &DbZnak{
		dbInfo: info,
	}
	if info == nil {
		return nil, fmt.Errorf("%s dbinfo is nil", modError)
	}
	// открываем сесиию в этом методе если нет ошибки
	if err := db.Check(); err != nil {
		return nil, fmt.Errorf("%s error check %v", modError, err)
	}
	return db, nil
}

func (c *DbZnak) Close() (err error) {
	return c.dbSession.Close()
}

func (c *DbZnak) Sess() db.Session {
	return c.dbSession
}

func (c *DbZnak) Version() int64 {
	return c.version
}

func (c *DbZnak) Info() dbscan.DbInfo {
	if c.dbInfo == nil {
		return dbscan.DbInfo{}
	}
	return *c.dbInfo
}

func (c *DbZnak) InfoType() dbscan.DbInfoType {
	return c.infoType
}
