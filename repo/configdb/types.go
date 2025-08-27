package configdb

import (
	_ "embed"
	"fmt"

	"github.com/mechiko/dbscan"
	"github.com/upper/db/v4"
	"go.uber.org/zap"
)

const modError = "selfdb"

type DbConfig struct {
	logger    *zap.SugaredLogger
	dbSession db.Session // открытый хэндл тут
	dbInfo    *dbscan.DbInfo
	infoType  dbscan.DbInfoType
	version   int64
}

func New(logger *zap.SugaredLogger, info *dbscan.DbInfo, infoType dbscan.DbInfoType) (*DbConfig, error) {
	db := &DbConfig{
		logger: logger,
		dbInfo: info,
	}
	if info == nil {
		return nil, fmt.Errorf("%s dbinfo is nil", modError)
	}
	// persist type for InfoType()
	db.infoType = infoType
	// открываем сесиию в этом методе если нет ошибки
	if err := db.Check(); err != nil {
		return nil, fmt.Errorf("%s error check %v", modError, err)
	}
	return db, nil
}

func (c *DbConfig) Close() error {
	if c.dbSession == nil {
		return nil
	}
	return c.dbSession.Close()
}

func (c *DbConfig) Sess() db.Session {
	return c.dbSession
}

func (c *DbConfig) Version() int64 {
	return c.version
}

func (c *DbConfig) Info() dbscan.DbInfo {
	if c.dbInfo == nil {
		return dbscan.DbInfo{}
	}
	return *c.dbInfo
}

func (c *DbConfig) InfoType() dbscan.DbInfoType {
	return c.infoType
}
