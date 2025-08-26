package zaplog

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

var Logger *zap.Logger

type ZapLog struct {
	logs map[LogName]*zap.Logger
}

func (z *ZapLog) GetLogger(name string) (*zap.Logger, error) {
	if isValidLogName(name) {
		if log, ok := z.logs[LogName(name)]; ok {
			return log, nil
		} else {
			return nil, fmt.Errorf("%s is not created", name)
		}
	} else {
		return nil, fmt.Errorf("%s not valid", name)
	}
}

func (z *ZapLog) Shutdown() {
	fmt.Println("zap log shutdown")
	for _, log := range z.logs {
		log.Sync()
	}
}

func New(outConfig map[string][]string, debug bool) (*ZapLog, error) {
	// проверяем мапу настройки логов
	for key, val := range outConfig {
		if !isValidLogName(key) {
			return nil, fmt.Errorf("wrong name %s", key)
		}
		if len(val) < 1 {
			return nil, fmt.Errorf("string array must be min lenght 1 elements")
		}
	}
	z := &ZapLog{
		logs: make(map[LogName]*zap.Logger),
	}
	err := z.init(outConfig, debug)
	if err != nil {
		return nil, fmt.Errorf("init zap logger %v", err)
	}
	return z, nil
}

//	var outConfig = map[string][]string{
//		"logger":   []string{"stdout"},
//		"echo":     []string{"stdout", "echo_name"},
//		"reductor": []string{"stderr"},
//		"true":     []string{"stdout", "true_name"},
//	}
func (z *ZapLog) Run(ctx context.Context) error {
	// ожидаем завершения контекста
	<-ctx.Done()
	fmt.Println("zaplog receive ctx shutdown")
	z.Shutdown()
	return nil
}

func (z *ZapLog) init(outConfig map[string][]string, debug bool) (err error) {
	for key, output := range outConfig {
		if isValidLogName(key) {
			lg, err := createLogger(output, debug)
			if err != nil {
				return fmt.Errorf("name %s %w", key, err)
			}
			z.logs[LogName(key)] = lg
		} else {
			return fmt.Errorf("wrong name for logger %s", key)
		}
	}
	// для совместимости основной логер пропишем в глобальную переменную
	if l, ok := z.logs[LogNameLogger]; ok {
		Logger = l
	}
	return nil
}
