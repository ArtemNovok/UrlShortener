package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"development"`
	StorageConn string     `yaml:"storage_conn" env-required:"true"`
	HttpServer  HttpServer `yaml:"http_server" env-required:"true"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	TimeOut     time.Duration `yaml:"timeout", env-default:"4s"`
	IdleTimeOut time.Duration `yaml:"idle_timeout", env-default:"60s"`
	Email       string        `yaml:"email" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true"`
}

func MustLoad() *Config {

	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	configPath := res
	if configPath == "" {
		log.Fatal("empty config path")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config path points to nonexistent file: %s", configPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("failed to read config file: %s", configPath)
	}
	return &cfg
}
