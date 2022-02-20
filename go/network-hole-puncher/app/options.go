package app

type Config struct {
	connMW []MW
}

type Option func(cfg *Config)

func newConfig(options ...Option) *Config {
	cfg := new(Config)
	for _, o := range options {
		o(cfg)
	}
	return cfg
}

func ConnOption(mw ...MW) Option {
	return func(cfg *Config) {
		cfg.connMW = append(cfg.connMW, mw...)
	}
}
