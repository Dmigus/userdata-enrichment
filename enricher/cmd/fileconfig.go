package main

import (
	"enricher/internal/app"
	"fmt"

	"github.com/vrischmann/envconfig"
)

func parseConfig() (*app.Config, error) {
	var conf app.Config
	err := envconfig.Init(&conf)
	if err != nil {
		return nil, fmt.Errorf("fatal error with config file: %w", err)
	}
	return &conf, nil
}
