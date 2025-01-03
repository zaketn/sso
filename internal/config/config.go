package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env-required:"true"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
}

type GRPCConfig struct {
	Port    int
	Timeout time.Duration
}

func MustLoad() *Config {
	cfgPath := fetchConfigPath()

	return MustLoadByPath(cfgPath)
}

func MustLoadByPath(cfgPath string) *Config {
	if cfgPath == "" {
		panic("config file is not specified")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("failed to read config file %s", err))
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		panic(fmt.Sprintf("failed to load config %s", err))
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config-path", "", "Path to the app config")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
