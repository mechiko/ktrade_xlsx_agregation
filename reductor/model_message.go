package reductor

import (
	"fmt"
)

type Model struct {
	Inn         string
	Kpp         string
	Order       int64
	Utilisation []int64
}

type ModelList map[string]Model

// формат сообщения
type Message struct {
	Sender string
	Page   string
	Model  Model
}

type IConfig interface {
	SetInConfig(key string, value interface{}, save ...bool) error
	GetKeyString(name string) string
	GetByName(name string) interface{}
}

func (m *Model) Sync(cfg IConfig) {
	// cfg.SetInConfig("entirely", m.Entirely, true)
	// cfg.SetInConfig("inn", m.Inn, true)
	// cfg.SetInConfig("ssccprefix", m.PrefixSSCC, true)
	// cfg.SetInConfig("ssccstartnumber", m.StartNumberSSCC, true)
	// cfg.SetInConfig("perpallet", m.PerPallet, true)
}

func (m *Model) Read(cfg IConfig) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	// val := cfg.GetByName("ssccstartnumber")
	// fmt.Printf("%T %v\n", val, val)
	// m.Entirely = cfg.GetByName("entirely").(bool)
	// m.Inn = cfg.GetKeyString("inn")
	// m.PrefixSSCC = cfg.GetKeyString("ssccprefix")
	// m.StartNumberSSCC = int(cfg.GetByName("ssccstartnumber").(int64))
	// m.PerPallet = int(cfg.GetByName("perpallet").(int64))
	return nil
}
