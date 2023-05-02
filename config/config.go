package config

import (
	"github.com/pkg/errors"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/source/env"
)

func Get() Config {
	return *cfg
}

func Load() error {
	config, err := config.NewConfig(config.WithSource(env.NewSource()))
	if err != nil {
		return errors.Wrap(err, "config.New")
	}
	if err := config.Load(); err != nil {
		return errors.Wrap(err, "config.Load")
	}
	if err := config.Scan(cfg); err != nil {
		return errors.Wrap(err, "config.Scan")
	}
	return nil
}
