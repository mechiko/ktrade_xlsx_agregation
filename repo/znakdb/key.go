package znakdb

import (
	"fmt"
)

func (c *DbZnak) Key(k string) (out string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic %v", r)
		}
	}()

	return "", nil
}
