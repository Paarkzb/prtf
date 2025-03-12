package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string     `yaml:"env" env-default:"local"`
	DB   DBConfig   `yaml:"db"`
	RDB  RDBConfig  `yaml:"rdb"`
	HTTP HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type RDBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); err != nil {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("uncorrect config file " + err.Error())
	}

	return &cfg
}
