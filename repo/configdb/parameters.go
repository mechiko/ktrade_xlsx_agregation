package configdb

import "github.com/upper/db/v4"

type Parameters struct {
	Name  string `db:"name"`
	Value string `db:"value"`
}

func (p *Parameters) Store(sess db.Session) db.Store {
	return sess.Collection("parameters")
}
