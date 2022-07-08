package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/fx"
)

var Module = fx.Module("config",
	fx.Provide(NewConfig),
)

type Config struct {
	TraceConfig *TraceConfig `env:",prefix=TRACE_"`
}

type TraceConfig struct {
	URL string `env:"URL,default=http://localhost:14268/api/traces"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
