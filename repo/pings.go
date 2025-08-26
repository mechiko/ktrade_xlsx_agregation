package repo

import "github.com/mechiko/dbscan"

func (r *Repository) Ping(t dbscan.DbInfoType) (result bool) {
	defer func() {
		if rec := recover(); rec != nil {
			result = false
		}
	}()
	info := r.dbs.Info(t)
	if info == nil {
		return false
	}
	if err := info.IsConnected(); err != nil {
		return false
	}
	return true
}
