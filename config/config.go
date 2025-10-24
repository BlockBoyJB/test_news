package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP HTTP
	Log  Log
	PG   PG
	JWT  JWT
}

type HTTP struct {
	Port string `env-required:"true" env:"HTTP_PORT"`
}

type Log struct {
	Level  string `env-required:"true" env:"LOG_LEVEL"`
	Output string `env-required:"true" env:"LOG_OUTPUT"`
}

type PG struct {
	Url string `env-required:"true" env:"PG_URL"`
}

type JWT struct {
	Key string `env-required:"true" env:"JWT_KEY"`
}

func NewConfig() (Config, error) {
	c := Config{}
	if err := cleanenv.ReadEnv(&c); err != nil {
		return Config{}, fmt.Errorf("error reading config env: %w", err)
	}
	return c, nil
}
