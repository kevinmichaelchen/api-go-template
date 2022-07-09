package grpc

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

func NewConfig() (cfg Config, err error) {
	err = envconfig.Process(context.Background(), &cfg)
	return
}

type Config struct {
	GRPCConfig    *NestedConfig  `env:",prefix=GRPC_"`
	ConnectConfig *ConnectConfig `env:",prefix=GRPC_CONNECT"`
}

type NestedConfig struct {
	Port int `env:"PORT,default=8080"`
}

type ConnectConfig struct {
	Host string `env:"HOST,default=localhost"`
	Port int    `env:"PORT,default=8081"`
}
