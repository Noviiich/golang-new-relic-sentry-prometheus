package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	StoragePath string `yaml:"storage_path"`
	Database    `yaml:"db" env-required:"true"`
	NewRelic    `yaml:"new_relic" env-required:"true"`
	Sentry      `yaml:"sentry" env-required:"true"`
	HTTPServer  `yaml:"http_server" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8082"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
}

type Database struct {
	Type                 string        `yaml:"type"`
	Host                 string        `yaml:"host"`
	Username             string        `yaml:"username"`
	Password             string        `yaml:"password"`
	Port                 string        `yaml:"port"`
	DBName               string        `yaml:"db_name"`
	SSLMode              string        `yaml:"ssl_mode"`
	MaxDBConnections     int           `yaml:"max_db_connections"`
	MaxDBIdleConnections int           `yaml:"max_db_idle_connections"`
	MaxDBLifetime        time.Duration `yaml:"max_db_lifetime"`
}

type NewRelic struct {
	AppName string `yaml:"app_name"`
	License string `yaml:"license"`
}

type Sentry struct {
	Dsn string `env:"SENTRY_DSN"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
