package config

import (
	"log"
	"os"
	"pg-sh-scripts/internal/config/api"
	"pg-sh-scripts/internal/config/postgres"
	"pg-sh-scripts/internal/config/project"
	"pg-sh-scripts/internal/config/server"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Project  project.Config
	Server   server.Config
	Api      api.Config      `yaml:"api"`
	Postgres postgres.Config `yaml:"postgres"`
}

var (
	cfgInstance *Config
	cfgOnce     sync.Once
)

func validateConfigPath(cfgPath string) error {
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return err
	}
	return nil
}

func setConfig(cfg *Config, cfgPath string) error {
	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		return err
	}
	return nil
}

func GetConfig() *Config {
	cfgOnce.Do(func() {
		var cfg Config

		cfgPath := "./config/app/main.yaml"

		if err := validateConfigPath(cfgPath); err != nil {
			log.Fatalf("Config path error: %v", err)
		}
		if err := setConfig(&cfg, cfgPath); err != nil {
			log.Fatalf("Config create error: %v", err)
		}

		cfgInstance = &cfg
	})

	return cfgInstance
}
