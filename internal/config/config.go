package config

import (
	"log"
	"path"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database `yaml:"database"`
	Stan     `yaml:"stan"`
}

type Database struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres""`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"db_name" env-default:"postgres"`
}

type Stan struct {
	ClusterID string `yaml:"cluster_id"`
	ClientID  string `yaml:"client_id"`
	SvrURL    string `yaml:"svr_url"`
	Subject   string `yaml:"subject"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(path.Join("config", "config.ex.yaml"), &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
