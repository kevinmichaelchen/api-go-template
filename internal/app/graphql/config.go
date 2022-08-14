package graphql

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

func NewConfig() (cfg Config, err error) {
	err = envconfig.Process(context.Background(), &cfg)
	return
}

type Config struct {
	GraphQLConfig *NestedConfig `env:",prefix=GRAPHQL_"`
}

type NestedConfig struct {
	Host string `env:"HOST,default=localhost"`
	Port int    `env:"PORT,default=8082"`
}
