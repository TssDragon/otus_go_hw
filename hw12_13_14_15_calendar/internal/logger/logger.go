package logger

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger   *zap.Logger
	logLevel zapcore.Level
}

func New(config config.LoggerConf) *Logger {
	var logLevel zapcore.Level
	switch config.Level {
	case "INFO":
		logLevel = zap.InfoLevel
	case "WARN":
		logLevel = zap.WarnLevel
	case "ERROR":
		logLevel = zap.ErrorLevel
	default:
		logLevel = zap.DebugLevel
	}

	fileName := time.Now().Format("2006_01_02") + "_log.log"
	filePath := path.Join(config.Dir, fileName)

	if _, err := os.Stat(config.Dir); os.IsNotExist(err) {
		if err = os.Mkdir(config.Dir, 0o766); err != nil {
			log.Fatal(err)
		}
	}

	l, err := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(logLevel),
		OutputPaths: []string{filePath},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: logLevel.String(),
		},
	}.Build()
	if err != nil {
		log.Fatal(err)
	}

	return &Logger{
		logger:   l,
		logLevel: logLevel,
	}
}

func (l Logger) Info(msg string) {
	l.logger.Info(msg)
	defer l.logger.Sync()
}

func (l Logger) Error(msg string) {
	l.logger.Error(msg)
	defer l.logger.Sync()
}

func (l Logger) Warn(msg string) {
	l.logger.Warn(msg)
	defer l.logger.Sync()
}
