package reductor

import (
	"sync"

	"go.uber.org/zap"
)

type Reductor struct {
	mutex  sync.Mutex
	in     chan Message
	logger *zap.SugaredLogger
	model  Model
}

var once sync.Once
var instance *Reductor

// по умолчанию страница пустая строка
// дальше придумаю как добавлять страницы методом
func New(model *Model, logger *zap.SugaredLogger) *Reductor {
	once.Do(func() {
		instance = &Reductor{
			in:     make(chan Message, 5),
			logger: logger,
		}
		instance.model = *model
	})
	return instance
}

func Instance() *Reductor {
	return instance
}

func (rdc *Reductor) Model(page string) (mdl Model) {
	rdc.mutex.Lock()
	defer rdc.mutex.Unlock()
	mdl = rdc.model
	return mdl
}

// записываем модель,  возвращаем предыдущую
func (rdc *Reductor) SetModel(page string, model *Model) {
	rdc.mutex.Lock()
	defer rdc.mutex.Unlock()
	rdc.model = *model
}

func (rdc *Reductor) Logger() *zap.SugaredLogger {
	return rdc.logger
}
