package main

import (
	"enricher/internal/app"
	"flag"
	"fmt"
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
	configName := viper.GetString(configNameFlag)
	viper.SetConfigFile(configName)
	err = viper.ReadInConfig() // Find and read the config file
	return err
}

func parseConfig() (*app.Config, error) {
	conf := app.Config{}
	err := viper.Unmarshal(&conf)
	if err != nil {
		return nil, fmt.Errorf("fatal error with config file: %w", err)
	}
	return &conf, nil
}
