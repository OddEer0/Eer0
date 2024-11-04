package eapp

import (
	"context"
)

type (
	Info struct {
		name, version, description string
	}

	UserConfigInterceptor = func(ctx context.Context, conf any) (any, error)
	LibConfigInterceptor  = func(ctx context.Context, conf *LibConfig) (*LibConfig, error)

	App struct {
		info               Info
		configs            *Configs
		err                error
		beforeHandlers     map[string]BeforeHandler
		afterHandlers      map[string]AfterHandler
		jobs               map[string]Job
		userCfgInterceptor UserConfigInterceptor
		libCfgInterceptor  LibConfigInterceptor
	}

	ReadApp struct {
		app *App
	}
)

func (r ReadApp) LibConfig() (LibConfig, error) {
	if r.app != nil && r.app.configs != nil && r.app.configs.Lib != nil {
		return LibConfig{}, ErrLibConfigToBeNull
	}
	return *r.app.configs.Lib, nil
}

func (r ReadApp) Config() (any, error) {
	if r.app != nil && r.app.configs != nil && r.app.configs.Client != nil {
		return nil, ErrUserConfigToBeNull
	}
	return *r.app.configs.Lib, nil
}

func (r ReadApp) Err() error {
	return r.app.err
}

func (r ReadApp) Info() Info {
	return r.app.info
}
