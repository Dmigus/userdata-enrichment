package main

import (
	"enricher/internal/app"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	configNameFlag       = "config"
	postgresPasswordFile = "postgresPasswordFile"
)

func initViper() error {
	flag.String(configNameFlag, "./configs/local.json", "path to config file for notifier")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return fmt.Errorf("fatal error binding flags: %w", err)
	}
	_ = viper.BindEnv(postgresPasswordFile, `POSTGRES_PASSWORD_FILE`)
	configName := viper.GetString(configNameFlag)
	viper.SetConfigFile(configName)
	err = viper.ReadInConfig() // Find and read the config file
	return err
}

func parseConfig() (*app.Config, error) {
	conf := app.Config{}
	postgresPwd, err := readSecretFromFile(viper.GetString(postgresPasswordFile))
	if err != nil {
		return nil, err
	}
	viper.Set(`Repository.Password`, postgresPwd)
	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, fmt.Errorf("fatal error with config file: %w", err)
	}
	return &conf, nil
}

func readSecretFromFile(addr string) (string, error) {
	cleaned := filepath.Clean(addr)
	dataBytes, err := os.ReadFile(cleaned)
	if err != nil {
		return "", err
	}
	return string(dataBytes), nil
}
