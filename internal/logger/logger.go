package logger

import (
	"L0/internal/config"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

const (
	debugLevel   = "debug"
	releaseLevel = "release"
)

var (
	ErrUnknowParam = errors.New("Unknow parameter")
)

/*
@brief - Инициализация объекта логгера
@param config - Объект конфиг файла
@return *zap.Logger - Объект логгера, при возниконовении ошибки возвращается nil
@return error - Ошибка в процессе инициализации объекта логгера, в случае успеха возвращается nil
*/
func NewLogger(config *config.Config) (*zap.Logger, error) {
	var (
		err    error
		logger *zap.Logger
	)

	switch config.LogLevel {
	case debugLevel:
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		logger.Info("Logger debug level is on.")
	case releaseLevel:
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
		logger.Info("Logger info level is on.")
	default:
		return nil, fmt.Errorf("%w %v", ErrUnknowParam, config.LogLevel)
	}

	return logger, err
}
