package config

type Config struct {
	Nats Nats
}

var cfg *Config = &Config{}
