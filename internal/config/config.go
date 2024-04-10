package config

import (
	"fmt"
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
	Project  project.Config  `yaml:"project"`
	Server   server.Config   `yaml:"server"`
	Api      api.Config      `yaml:"api"`
	Postgres postgres.Config `yaml:"postgres"`
}

func getConfigPath() string { return "./config/app/main.yaml" }

func validateConfigPath(configPath *string) error {
	if *configPath == "" {
		return fmt.Errorf("empty config path")
	}
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		return err
	}
	return nil
}

func setConfig(config *Config, configPath *string) error {
	if err := cleanenv.ReadConfig(*configPath, config); err != nil {
		return err
	}
	return nil
}

func GetConfig() *Config {
	var (
		config Config
		once   sync.Once
	)

	once.Do(func() {
		configPath := getConfigPath()
		if err := validateConfigPath(&configPath); err != nil {
			log.Fatalf("Config path error: %v", err)
		}
		if err := setConfig(&config, &configPath); err != nil {
			log.Fatalf("Config create error: %v", err)
		}
	})

	return &config
}
