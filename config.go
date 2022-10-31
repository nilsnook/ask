package main

import (
	"path"

	"github.com/spf13/viper"
)

type config struct {
	Keys []string `mapstructure:"keys"`
}

func (c *config) setDefaults() {
	defaultKeys := []string{
		"~/.ssh/id_ed25519",
		"~/.ssh/id_rsa",
	}
	viper.SetDefault("keys", defaultKeys)
}

func (c *config) readConfigFileIn(configdir string) error {
	// set config name
	viper.SetConfigName("config")
	// set config type
	viper.SetConfigType("yaml")
	// add config path
	viper.AddConfigPath(configdir)
	// read from config file
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func (c *config) writeConfigFileIn(configdir string) error {
	// create configdir if not exist
	err := createDirIfNotExists(configdir)
	if err != nil {
		return err
	}
	// write config file
	configFile := path.Join(configdir, "config.yaml")
	if err = viper.WriteConfigAs(configFile); err != nil {
		return err
	}
	return nil
}
