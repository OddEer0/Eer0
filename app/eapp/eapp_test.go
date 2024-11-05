package eapp

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type clientConf struct {
	Message string `yaml:"message"`
}

type jobI struct {
	init, run, stop *int
}

func (j *jobI) Init(ctx context.Context, config any) error {
	*j.init = 1
	return nil
}

func (j *jobI) Run(ctx context.Context) error {
	*j.run = 1
	return nil
}

func (j *jobI) Close(ctx context.Context) error {
	*j.stop = 1
	return nil
}

func UserInter(ctx context.Context, cfg any) (any, error) {
	conf, ok := cfg.(*clientConf)
	if !ok {
		return nil, errors.New("invalid config")
	}
	conf.Message = "interceptor"
	return conf, nil
}

func LibInter(ctx context.Context, cfg *LibConfig) (*LibConfig, error) {
	cfg.Info["test"] = "interceptor"
	return cfg, nil
}

func TestApp(t *testing.T) {
	initOpt := &InitOptions{
		Name:        "my-app",
		Version:     "0.0.1",
		Description: "This is my app",
	}

	libCfg := &LibConfig{
		StartTimeout: 5 * time.Second,
		StopTimeout:  5 * time.Second,
	}

	t.Run("Should correct initialization", func(t *testing.T) {
		a := Init(initOpt)
		assert.Equal(t, a.info, Info{
			name:        initOpt.Name,
			version:     initOpt.Version,
			description: initOpt.Description,
		})
		assert.Equal(t, a.configs, &Configs{})
		assert.NotNil(t, a.beforeHandlers)
		assert.NotNil(t, a.afterHandlers)
		assert.NotNil(t, a.jobs)
		assert.Nil(t, a.err)
	})

	t.Run("Should return error with empty init opts", func(t *testing.T) {
		a := Init(&InitOptions{
			Name:        "my-app",
			Version:     "0.0.1",
			Description: "",
		})
		assert.Error(t, a.err)

		a = Init(&InitOptions{
			Name:        "my-app",
			Version:     "",
			Description: "my app",
		})
		assert.Error(t, a.err)

		a = Init(&InitOptions{
			Name:        "",
			Version:     "0.0.1",
			Description: "my app",
		})
		assert.Error(t, a.err)
	})

	t.Run("Should correct lib config", func(t *testing.T) {
		app := Init(initOpt).LibConfig(libCfg)
		assert.Equal(t, app.configs.Lib, libCfg)

		app = Init(&InitOptions{}).LibConfig(libCfg)
		assert.Nil(t, app.configs)
	})

	t.Run("Should correct user config parse", func(t *testing.T) {
		clCfg := &clientConf{
			Message: "testing message",
		}
		app := Init(initOpt).LibConfig(libCfg).WithConfig(clCfg)
		assert.NoError(t, app.err)
		assert.NotNil(t, app.configs.Client)
		cfg, ok := app.configs.Client.(*clientConf)
		assert.True(t, ok)
		assert.Equal(t, cfg.Message, "testing message")
	})

	t.Run("Should correct options", func(t *testing.T) {
		cfg := &clientConf{}
		opt := &InitOptions{
			Name:               "my-app",
			Version:            "0.0.1",
			Description:        "This is my app",
			UserCfgInterceptor: UserInter,
			LibCfgInterceptor:  LibInter,
		}
		app := Init(opt).
			LibConfig(libCfg).
			WithConfig(cfg)

		assert.Equal(t, cfg.Message, "interceptor")
		assert.Equal(t, app.Configs().Lib.Info["test"], "interceptor")
	})

	t.Run("Should correct after, before, job", func(t *testing.T) {
		bMsg := 0
		aMsg := 0
		jobInitMsg := 0
		jobStartMsg := 0
		jobCloseMsg := 0

		bHandle := BeforeHandler{
			Key: "1",
			Handler: func(ctx context.Context, cfg any) error {
				bMsg++
				return nil
			},
		}
		aHandle := AfterHandler{
			Key: "1",
			Handler: func(ctx context.Context, cfg any) error {
				aMsg++
				return nil
			},
		}

		err := Init(initOpt).
			LibConfig(libCfg).
			BeforeHandle(bHandle).
			AfterHandle(aHandle).
			WithJobs(JobOption{
				Key: "1",
				Job: &jobI{
					init: &jobInitMsg,
					run:  &jobStartMsg,
					stop: &jobCloseMsg,
				},
			}).
			Start()

		assert.NoError(t, err)
		assert.Equal(t, 1, bMsg)
		assert.Equal(t, 1, aMsg)
		assert.Equal(t, 1, jobInitMsg)
		assert.Equal(t, 1, jobStartMsg)
		assert.Equal(t, 1, jobCloseMsg)
	})
}
