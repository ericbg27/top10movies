package logger

import (
	"fmt"
	"strings"

	"github.com/ericbg27/top10movies-api/src/utils/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Print(v ...interface{})
}

type logger struct {
	log *zap.Logger
}

var (
	log logger
)

func init() {
	cfg := config.GetConfig()

	logConfig := zap.Config{
		OutputPaths: []string{getOutput(cfg)},
		Level:       zap.NewAtomicLevelAt(getLevel(cfg)),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error
	if log.log, err = logConfig.Build(); err != nil {
		panic(err)
	}
}

func getLevel(cfg config.Config) zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(cfg.Logger.LogLevel)) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func getOutput(cfg config.Config) string {
	output := strings.TrimSpace(cfg.Logger.LogOutput)

	if output == "" {
		return "stdout"
	}

	return output
}

func GetLogger() Logger {
	return log
}

func (l logger) Print(v ...interface{}) {
	Info(fmt.Sprintf("%v", v))
}

func Info(msg string, tags ...zap.Field) {
	log.log.Info(msg, tags...)
	log.log.Sync()
}

func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.log.Error(msg, tags...)
	log.log.Sync()
}
