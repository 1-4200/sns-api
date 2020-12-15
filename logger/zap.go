package logger

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sns-api/config"
	"strings"
)

func NewZapLogger(c *config.Config) (*zap.SugaredLogger, error) {
	var cfg zap.Config

	switch strings.ToLower(c.Logger.Environment) {
	case "dev", "development":
		cfg = zap.NewDevelopmentConfig()
	case "prod", "production":
		cfg = zap.NewProductionConfig()
	default:
		return nil, errors.New("logger environment not supported")
	}

	cfg.Level = zap.NewAtomicLevelAt(getLevel(c.Logger.LogLevel))
	cfg.OutputPaths = []string{c.Logger.FileName}
	log, err := cfg.Build()
	if err != nil {
		return nil, errors.New("zap logger build constructs failed")
	}
	return log.Sugar(), nil
}

func getLevel(level string) zapcore.Level {
	var zapLevel zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		zapLevel = zap.DebugLevel

	case "info":
		zapLevel = zap.InfoLevel

	case "warn", "warning":
		zapLevel = zap.WarnLevel

	case "error":
		zapLevel = zap.ErrorLevel
	}
	return zapLevel
}

