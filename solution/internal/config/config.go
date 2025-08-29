package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"dev" env-required:"true"`
	HTTPServer `yaml:"http_server" env-required:"true"`
	DbUser     string `yaml:"db_user" env:"DB_USER" env-required:"true"`
	DbPassword string `yaml:"db_password" env:"DB_PASSWORD" env-required:"true"`
	DbName     string `yaml:"db_name" env:"DB_NAME" env-required:"true"`
	DbHost     string `yaml:"db_host" env:"DB_HOST" env-required:"true"`
	DbPort     string `yaml:"db_port" env:"DB_PORT" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env:"ADDRESS" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"30s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	return &cfg
}
