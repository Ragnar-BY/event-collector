package config

import (
	"errors"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is config struct
type Config struct {
	ClickhouseAddr     string `env:"CLICKHOUSE_ADDRESS"`
	ClickhouseUser     string `env:"CLICKHOUSE_USER"`
	ClickhousePassword string `env:"CLICKHOUSE_PASSWORD"`
	ClickhouseDatabase string `env:"CLICKHOUSE_DATABASE"`

	ServerAddress string `env:"SERVER_ADDRESS"`

	NumberOfInsertThreads int `env:"NUMBER_OF_INSERT_THREADS"`
	ChannelCapacity       int `env:"CHANNEL_CAPACITY"`
}

// LoadConfig loads config from .env file
func LoadConfig(path string) (config Config, err error) {
	err = cleanenv.ReadConfig(path, &config)
	if errors.Is(err, os.ErrNotExist) {
		err = cleanenv.ReadEnv(&config)
	}
	return
}
