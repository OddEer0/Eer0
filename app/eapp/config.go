package eapp

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

type (
	Configs struct {
		Client any
		Lib    *LibConfig
	}

	LibConfig struct {
		StartTimeout time.Duration
		StopTimeout  time.Duration
		Info         map[string]any
	}
)

func (a *App) LibConfig(conf *LibConfig) *App {
	if a.err != nil {
		return a
	}

	if conf.Info == nil {
		conf.Info = make(map[string]any)
	}

	if a.libCfgInterceptor != nil {
		var err error
		conf, err = a.libCfgInterceptor(AppWithContext(context.Background(), a), conf)
		if err != nil {
			a.err = err
			return a
		}
	}

	a.configs.Lib = conf

	return a
}

func (a *App) Configs() *Configs {
	return a.configs
}

func (a *App) WithConfig(cfg any) *App {
	if a.Err() != nil {
		return a
	}

	if a.userCfgInterceptor != nil {
		var err error
		cfg, err = a.userCfgInterceptor(AppWithContext(context.Background(), a), cfg)
		if err != nil {
			a.err = errors.Wrap(err, "[App] userCfgInterceptor")
			return a
		}
	}

	a.configs.Client = cfg
	return a
}
