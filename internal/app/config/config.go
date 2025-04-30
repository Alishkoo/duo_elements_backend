package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/mcuadros/go-defaults"
)

type Config struct {
	HTTPServer HTTPServerConfig `envPrefix:"HTTP_"`
	DB         DBConfig         `envPrefix:"DB_"`
}

type HTTPServerConfig struct {
	Host string `env:"HOST" default:"localhost"`
	Port string `env:"PORT" default:"8080"`
}

type DBConfig struct {
	Host     string `env:"HOST" default:"localhost"`
	Port     string `env:"PORT" default:"5432"`
	User     string `env:"USER" default:"user"`
	Password string `env:"PASSWORD" default:"password"`
	Name     string `env:"NAME" default:"mydb"`
}

func NewConfig(filenames ...string) (*Config, error) {
	for _, file := range filenames {
		if err := godotenv.Load(file); err != nil {
			log.Printf("Не удалось загрузить .env файл %s: %v", file, err)
		}
	}

	cfg := &Config{}
	defaults.SetDefaults(cfg)

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
