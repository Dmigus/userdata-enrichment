package main

import (
	"errors"
	"github.com/vrischmann/envconfig"

	"github.com/samber/lo"
)

func parseConfigs(configs ...any) error {
	errs := lo.Map(configs, func(conf any, _ int) error {
		return envconfig.Init(conf)
	})
	return errors.Join(errs...)
}
