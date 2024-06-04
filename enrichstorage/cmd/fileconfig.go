package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
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

func parseConfigs(configs ...any) error {
	postgresPwd, err := readSecretFromFile(viper.GetString(postgresPasswordFile))
	if err != nil {
		return fmt.Errorf("error reading postgres password: %w", err)
	}
	viper.Set(`Repository.Password`, postgresPwd)
	errs := lo.Map(configs, func(conf any, _ int) error {
		return viper.Unmarshal(conf)
	})
	return errors.Join(errs...)
}

func readSecretFromFile(addr string) (string, error) {
	cleaned := filepath.Clean(addr)
	dataBytes, err := os.ReadFile(cleaned)
	if err != nil {
		return "", err
	}
	return string(dataBytes), nil
}
