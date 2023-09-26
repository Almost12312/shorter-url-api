package config

import "time"

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`

	HttpServer `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"storage_path" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"timeout" env-default:"60s"`
}

func MustLoad() {

}
