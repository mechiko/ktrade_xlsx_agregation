package repo

import "github.com/mechiko/dbscan"

func (r *Repository) Is(t dbscan.DbInfoType) bool {
	if di := r.dbs.Info(t); di != nil {
		return true
	}
	return false
}
