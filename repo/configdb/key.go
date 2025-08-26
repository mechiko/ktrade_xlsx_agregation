package configdb

import (
	"fmt"

	"github.com/upper/db/v4"
)

func (c *DbConfig) Key(k string) (out string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic %v", r)
		}
	}()

	param := &Parameters{}
	if err := c.dbSession.Get(param, db.Cond{"name": k}); err != nil {
		return "", fmt.Errorf("%s %v", modError, err)
	}
	return param.Value, nil
}
