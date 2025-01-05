package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string `yaml:"env" env-default:"local"`
	HTTPServer    `yaml:"http_server"`
	PgConfig      `yaml:"pg_config"`
	MetricsConfig `yaml:"metrics_config"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
	Session     Session       `yaml:"session"`
}

type Session struct {
	SecretKey string `yaml:"secret_key" env-default:"secret_key"`
	Secure    bool   `yaml:"secure" env-default:"false"`
	MaxAge    int    `yaml:"max_age" env-default:"604800"`
}
type PgConfig struct {
	Host          string `yaml:"host" env-default:"localhost"`
	Port          int    `yaml:"port" env-default:"5432"`
	User          string `yaml:"user" env-default:"todo_list"`
	Password      string `yaml:"password" env-default:"pg"`
	DbName        string `yaml:"db_name" env-default:"todo_list"`
	MigrationsDir string `yaml:"migrations_dir" env-default:"./migrations"`
}

type MetricsConfig struct {
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
