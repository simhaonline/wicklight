package config

import (
	"wicklight/logger"

	"github.com/BurntSushi/toml"
)

// Conf default default
var Conf Config

// ReadConfig read config from file
func ReadConfig(configFile string) {
	if _, err := toml.DecodeFile(configFile, &Conf); err != nil {
		logger.Fatal("[config] failed to read config:", configFile, err)
	}

	if Conf.Log.Level != 0 {
		logger.SetLevel(Conf.Log.Level)
	}
	if Conf.Log.File != "" {
		logger.SetOutput(Conf.Log.File)
	}
	logger.Debug("[config] read config successfully:", Conf)
}
