package logger

import (
	"fmt"
	"strings"

	"github.com/ericbg27/top10movies-api/src/utils/config"
	"github.com/jackc/pgx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Print(v ...interface{})
	Log(level pgx.LogLevel, msg string, data map[string]interface{})
}

type logger struct {
	log *zap.Logger
}

var (
	log logger
)

const (
	encodingString = "json"
	levelKey       = "level"
	timeKey        = "time"
	messageKey     = "message"
)

func init() {
	cfg := config.GetConfig()

	var err error
	log.log, err = setupLogger(cfg)
	if err != nil {
		panic(err)
	}
}

func setupLogger(c *config.Config) (*zap.Logger, error) {
	logConfig := zap.Config{
		OutputPaths: []string{getOutput(*c)},
		Level:       zap.NewAtomicLevelAt(getLevel(*c)),
		Encoding:    encodingString,
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     levelKey,
			TimeKey:      timeKey,
			MessageKey:   messageKey,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	l, err := logConfig.Build()
	if err != nil {
		return nil, err
	}

	return l, nil
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

func (l logger) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
	fields := make([]zap.Field, len(data))
	i := 0
	for k, v := range data {
		fields[i] = zap.Reflect(k, v)
		i++
	}

	l.log.Info(msg, fields...)
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
