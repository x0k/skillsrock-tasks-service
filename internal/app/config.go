package app

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type LoggerConfig struct {
	Level       string `env:"LOGGER_LEVEL" env-default:"info"`
	HandlerType string `env:"LOGGER_HANDLER_TYPE" env-default:"text"`
}

type PgConfig struct {
	ConnectionURI string `env:"PG_CONNECTION_URI" env-required:"true"`
	MigrationsURI string `env:"PG_MIGRATIONS_URI" env-default:"file://db/migrations"`
}

type RedisConfig struct {
	ConnectionURI string `env:"REDIS_CONNECTION_URI" env-required:"true"`
}

type ServerConfig struct {
	Address string `yaml:"address" env:"SERVER_ADDRESS" env-default:"0.0.0.0:8080"`
}

type AuthConfig struct {
	Secret        string        `yaml:"secret" env:"AUTH_SECRET" env-required:"true"`
	TokenLifetime time.Duration `yaml:"token_lifetime" env:"AUTH_TOKEN_LIFETIME" env-default:"10m"`
}

type Config struct {
	Logger   LoggerConfig `yaml:"logger"`
	Postgres PgConfig     `yaml:"postgres"`
	Redis    RedisConfig  `yaml:"redis"`
	Server   ServerConfig `yaml:"server"`
	Auth     AuthConfig   `yaml:"auth"`
}

func MustLoadConfig(configPath string) *Config {
	cfg := &Config{}
	var cfgErr error
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfgErr = cleanenv.ReadEnv(cfg)
	} else if err == nil {
		cfgErr = cleanenv.ReadConfig(configPath, cfg)
	} else {
		cfgErr = err
	}
	if cfgErr != nil {
		log.Fatalf("cannot read config: %s", cfgErr)
	}
	return cfg
}
