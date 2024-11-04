package ecosystem

import (
	"context"
	"github.com/OddEer0/Eer0/app/eapp"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"os"
	"time"
)

func AdapterParseConfigFromYaml(path string) eapp.UserConfigInterceptor {
	return func(ctx context.Context, config any) (any, error) {
		_, err := os.Stat(path)
		if err != nil {
			return nil, errors.Wrap(err, "cannot stat config file")
		}

		err = cleanenv.ReadConfig(path, config)
		if err != nil {
			return nil, errors.Wrap(err, "cannot read config file")
		}

		return config, nil
	}
}

type LibCfg struct {
	StartTimeout time.Duration  `yaml:"start_timeout"`
	StopTimeout  time.Duration  `yaml:"stop_timeout"`
	Info         map[string]any `yaml:"info,omitempty"`
}

func AdapterParseLibConfigFromYaml(path string) eapp.LibConfigInterceptor {
	return func(ctx context.Context, conf *eapp.LibConfig) (*eapp.LibConfig, error) {
		_, err := os.Stat(path)
		if err != nil {
			return nil, errors.Wrap(err, "cannot stat config file")
		}

		lCfg := &LibCfg{}
		err = cleanenv.ReadConfig(path, lCfg)
		if err != nil {
			return nil, errors.Wrap(err, "cannot read config file")
		}

		conf.Info = lCfg.Info
		conf.StartTimeout = lCfg.StartTimeout
		conf.StopTimeout = lCfg.StopTimeout

		return conf, nil
	}
}
