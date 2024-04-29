package postgres

import "time"

type Config struct {
	Database          string        `yaml:"database"          env:"POSTGRES_DB"       env-required:"true"`
	Username          string        `yaml:"username"          env:"POSTGRES_USER"     env-required:"true"`
	Password          string        `yaml:"password"          env:"POSTGRES_PASSWORD" env-required:"true"`
	Host              string        `yaml:"host"              env:"POSTGRES_HOST"     env-required:"true"`
	Port              string        `yaml:"port"              env:"POSTGRES_PORT"     env-required:"true"`
	RetryCount        int           `yaml:"retryCount"`
	RetrySleepSeconds time.Duration `yaml:"retrySleepSeconds"`
}
