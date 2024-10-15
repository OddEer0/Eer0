package eapp

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"os"
	"time"
)

type (
	Configs struct {
		Client any
		Lib    *LibConfig
	}

	LibConfig struct {
		StartTimeout     time.Duration `yaml:"start_timeout"`
		StopTimeout      time.Duration `yaml:"stop_timeout"`
		ClientConfigPath string        `yaml:"client_config_path"`
	}
)

func (a *app) LibConfig(conf *LibConfig) App {
	if a.err != nil {
		return a
	}

	a.configs.Lib = conf

	return a
}

func (a *app) LibConfigByYaml(filePath string) App {
	if a.err != nil {
		return a
	}

	cfg := &LibConfig{}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		a.err = errors.Wrap(err, "[app.LibConfigByYaml] os.Stat")
		return a
	}

	err := cleanenv.ReadConfig(filePath, cfg)
	if err != nil {
		a.err = errors.Wrap(err, "[app.LibConfigByYaml] cleanenv.ReadConfig")
		return a
	}

	return a
}

func (a *app) WithConfig(cfg any) App {
	if a.Err() != nil {
		return a
	}

	path := a.configs.Lib.ClientConfigPath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		a.err = errors.Wrap(err, "[app.WithConfig] os.Stat")
		return a
	}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		a.err = errors.Wrap(err, "[app.WithConfig] cleanenv.ReadConfig")
		return a
	}

	a.configs.Client = cfg
	return a
}
