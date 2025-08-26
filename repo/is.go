package repo

import "github.com/mechiko/dbscan"

func (r *Repository) Is(t dbscan.DbInfoType) bool {
	if r == nil || r.dbs == nil {
		return false
	}
	return r.dbs.Info(t) != nil
}
