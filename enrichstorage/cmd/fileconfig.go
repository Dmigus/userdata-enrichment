package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/samber/lo"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	configNameFlag = "config"
)

func initViper() error {
	flag.String(configNameFlag, "./configs/local.json", "path to config file for notifier")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return fmt.Errorf("fatal error binding flags: %w", err)
	}
	_ = viper.BindEnv(`Storage.Password`, `POSTGRES_PASSWORD`)
	configName := viper.GetString(configNameFlag)
	viper.SetConfigFile(configName)
	err = viper.ReadInConfig() // Find and read the config file
	return err
}

func parseConfigs(configs ...any) error {
	errs := lo.Map(configs, func(conf any, _ int) error {
		return viper.Unmarshal(conf)
	})
	return errors.Join(errs...)
}
