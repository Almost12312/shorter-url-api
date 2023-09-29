package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`

	HttpServer `yaml:"http_server" env:"HTTP_SERVER"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle-timeout" env-default:"60s"`

	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true" env:"PASSWORD"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("Config path is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file is not found: %s", configPath)
	}

	var conf Config

	if err := cleanenv.ReadConfig(configPath, &conf); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	return &conf
}
