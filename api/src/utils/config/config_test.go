package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	testCfgName    = "configTest"
	testBadCfgName = "badConfigTest"
	testCfgType    = "yaml"
	testCfgPath    = ""
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestSetUpConfigSuccess(t *testing.T) {
	viper.Reset()

	var testCfg *Config
	var err error

	testCfg, err = setupConfig(testCfgName, testCfgPath, testCfgType)

	assert.NotNil(t, testCfg)
	assert.Nil(t, err)

	assert.EqualValues(t, "localhost", testCfg.Server.Host)
	assert.EqualValues(t, "8080", testCfg.Server.Port)

	assert.EqualValues(t, "Info", testCfg.Logger.LogLevel)
	assert.EqualValues(t, "", testCfg.Logger.LogOutput)

	assert.EqualValues(t, "localhost", testCfg.Database.Host)
	assert.EqualValues(t, 5000, testCfg.Database.Port)
	assert.EqualValues(t, "eric", testCfg.Database.User)
	assert.EqualValues(t, "1234", testCfg.Database.Password)
	assert.EqualValues(t, "dbtest", testCfg.Database.DbName)
	assert.EqualValues(t, "info", testCfg.Database.LogLevel)
}

func TestSetUpConfigFailureNoFile(t *testing.T) {
	viper.Reset()

	var testCfg *Config
	var err error

	testCfg, err = setupConfig("inexistentConfig", configType, configPath)

	assert.Nil(t, testCfg)
	assert.NotNil(t, err)
}

func TestSetUpConfigFailureUnmarshal(t *testing.T) {
	viper.Reset()

	var testCfg *Config
	var err error

	testCfg, err = setupConfig(testBadCfgName, configType, configPath)

	assert.Nil(t, testCfg)
	assert.NotNil(t, err)
}
