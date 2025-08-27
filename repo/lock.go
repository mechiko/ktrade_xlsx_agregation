package repo

import (
	"agregat/domain"
	"agregat/repo/configdb"
	"agregat/repo/znakdb"
	"fmt"

	"github.com/mechiko/dbscan"
)

// after Self() must be SelfClose() or deadlock
func (r *Repository) Lock(t dbscan.DbInfoType) (domain.RepoDB, error) {
	r.logger.Infof("repo Lock %v", t)
	info := r.dbs.Info(t)
	if info == nil {
		return nil, nil
	}
	mu, ok := r.dbMutex[t]
	if ok {
		mu.mutex.Lock()
	} else {
		return nil, fmt.Errorf("repo lock not present mutex %v", t)
	}
	switch t {
	case dbscan.Config:
		return configdb.New(r.logger, info, t)
	case dbscan.TrueZnak:
		return znakdb.New(info, t)
	default:
		mu.mutex.Unlock()
		return nil, fmt.Errorf("repo lock not present type mutex %v", t)
	}
}

func (r *Repository) Unlock(t dbscan.DbInfoType) error {
	r.logger.Infof("repo UnLock %v", t)
	mu, ok := r.dbMutex[t]
	if ok {
		mu.mutex.Unlock()
	} else {
		return fmt.Errorf("%s unlock not present mutex %v", modError, dbscan.Other)
	}
	return nil
}
