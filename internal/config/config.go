package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	HTTPServer `yaml:"http_server"`
	DBServer   `yaml:"db_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"http_address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DBServer struct {
	Address  string `yaml:"db_address" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"myuser"`
	Password string `yaml:"password" env-default:"mypassword"`
	DBname   string `yaml:"dbname" env-default:"mydatabase"`
}

// MustLoad loads the configuration from the specified path and returns it.
// It panics if the configuration cannot be loaded.
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH") // TODO: implement flag `--path`

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
