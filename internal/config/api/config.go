package api

type Config struct {
	Prefix         string   `env:"API_PREFIX" env-required:"true"`
	TrustedProxies []string `yaml:"trustedProxies"`
}
