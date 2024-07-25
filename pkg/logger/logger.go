package logger

import "go.uber.org/zap"

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger() (*Logger, error) {
	log, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &Logger{log.Sugar()}, nil
}
