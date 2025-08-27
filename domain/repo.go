package domain

import (
	"context"

	"github.com/mechiko/dbscan"
	"github.com/upper/db/v4"
)

type Repo interface {
	Is(dbscan.DbInfoType) bool
	Ping(dbscan.DbInfoType) bool
	Shutdown()
	Run(context.Context) error
	Lock(dbscan.DbInfoType) (RepoDB, error)
	Unlock(dbscan.DbInfoType) error
	Info(t dbscan.DbInfoType) *dbscan.DbInfo
	ListDbs() (out []dbscan.DbInfoType)
}

type RepoDB interface {
	Close() error
	Sess() db.Session
	Version() int64
	Info() dbscan.DbInfo
	InfoType() dbscan.DbInfoType
}
