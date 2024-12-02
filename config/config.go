package config

import (
	"errors"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	ENV            string         `env:"ENV" envDefault:"dev" mapstructure:"ENV"`
	PORT           string         `env:"PORT" envDefault:"8080" mapstructure:"PORT"`
	PostgresConfig PostgresConfig `envPrefix:"POSTGRES_" mapstructure:"POSTGRES"`
	JWT            JWTConfig      `envPrefix:"JWT_" mapstructure:"JWT"`
}

type JWTConfig struct {
	SecretKey string `env:"SECRET_KEY" envDefault:"secret" mapstructure:"SECRET_KEY"`
}
	
type PostgresConfig struct {
	Host     string `env:"HOST" envDefault:"localhost" mapstructure:"HOST"`
	Port     string `env:"PORT" envDefault:"5432" mapstructure:"PORT"`
	User     string `env:"USER" envDefault:"postgres" mapstructure:"USER"`
	Password string `env:"PASSWORD" envDefault:"postgres" mapstructure:"PASSWORD"`
	Database string `env:"DATABASE" envDefault:"postgres" mapstructure:"DATABASE"`
}

func NewConfig(envPath string) (*Config, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		return nil, errors.New("failed to load env")
	}

	cfg := new(Config)
	err = env.Parse(cfg)
	if err != nil {
		return nil, errors.New("failed to parse env")
	}
	return cfg, nil
}
