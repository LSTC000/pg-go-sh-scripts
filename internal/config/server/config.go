package server

type Config struct {
	Port string `env:"SERVER_PORT" env-required:"true"`
}
