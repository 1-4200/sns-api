package logger

import (
	"errors"
	"log"
	"sns-api/config"

	"go.uber.org/zap"
)

type Level interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
}

type Format interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

type Logging interface {
	Level
	Format
}

type Logger struct {
	ZapSugarLogger *zap.SugaredLogger
}

func NewLogger(c *config.Config) (Logging, error) {
	if c.Logger.Use == "zap" {
		z, er := NewZapLogger(c)
		if er != nil {
			log.Fatalf("can't initialize zap logger: %v", er)
			return nil, er
		}
		return &Logger{ZapSugarLogger: z}, nil
	}
	return nil, errors.New("logger not supported : " + c.Logger.Use)
}

func (l *Logger) Debug(args ...interface{}) {
	l.ZapSugarLogger.Debug(args)
}

func (l *Logger) Info(args ...interface{}) {
	l.ZapSugarLogger.Info(args)
}

func (l *Logger) Warn(args ...interface{}) {
	l.ZapSugarLogger.Warn(args)
}

func (l *Logger) Error(args ...interface{}) {
	l.ZapSugarLogger.Error(args)
}

func (l *Logger) Panic(args ...interface{}) {
	l.ZapSugarLogger.Panic(args)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.ZapSugarLogger.Fatal(args)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.ZapSugarLogger.Debugf(template, args)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.ZapSugarLogger.Infof(template, args)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.ZapSugarLogger.Warnf(template, args)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.ZapSugarLogger.Errorf(template, args)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.ZapSugarLogger.Panicf(template, args)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.ZapSugarLogger.Fatalf(template, args)
}
