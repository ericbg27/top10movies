package logger

import (
	"os"
	"testing"

	"github.com/ericbg27/top10movies-api/src/utils/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestSetupLoggerSuccess(t *testing.T) {
	testCfg := &config.Config{
		Logger: config.LoggerCfg{
			LogLevel:  "Info",
			LogOutput: "",
		},
	}

	l, err := setupLogger(testCfg)

	assert.NotNil(t, l)
	assert.Nil(t, err)

	assert.NotNil(t, l.Check(zap.InfoLevel, "Info"))
}

func TestSetupLoggerFailureWrongOutputPath(t *testing.T) {
	testCfg := &config.Config{
		Logger: config.LoggerCfg{
			LogLevel:  "Info",
			LogOutput: "file:///d:/a/wronglog.log",
		},
	}

	l, err := setupLogger(testCfg)

	assert.Nil(t, l)
	assert.NotNil(t, err)
}

func TestGetLevel(t *testing.T) {
	testCfg := config.Config{}

	testCfg.Logger.LogLevel = "Debug"
	assert.EqualValues(t, zap.DebugLevel, getLevel(testCfg))

	testCfg.Logger.LogLevel = "Info"
	assert.EqualValues(t, zap.InfoLevel, getLevel(testCfg))

	testCfg.Logger.LogLevel = "Error"
	assert.EqualValues(t, zap.ErrorLevel, getLevel(testCfg))

	testCfg.Logger.LogLevel = ""
	assert.EqualValues(t, zap.InfoLevel, getLevel(testCfg))
}

func TestGetOutput(t *testing.T) {
	testCfg := config.Config{}

	testCfg.Logger.LogOutput = ""
	assert.EqualValues(t, "stdout", getOutput(testCfg))

	testCfg.Logger.LogOutput = "     stdout    "
	assert.EqualValues(t, "stdout", getOutput(testCfg))
}
