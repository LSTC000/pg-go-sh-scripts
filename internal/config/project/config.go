package project

type Config struct {
	Mode string `env:"PROJECT_MODE" env-required:"true" env-description:"local/dev/prod"`
}
