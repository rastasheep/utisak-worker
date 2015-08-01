package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

var configPath = flag.String("config", "config/config.json", "Path to configuration file")

type Config struct {
	LogTo            string
	LogLevel         string
	ReadabilityToken string
	FeedRegistryPath string
	RedisDomain      string
	Redis            struct {
		Domain   string
		Database string
		Pool     string
		Process  string
	}
	Postgres struct {
		Host     string
		Port     string
		Database string
		Username string
		Password string
	}
}

func LoadConfig() *Config {
	flag.Parse()

	config := Config{}
	raw, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Could not read config file: %v", err))
	}

	if err = json.Unmarshal(raw, &config); err != nil {
		panic(fmt.Sprintf("Could not unmarshal config file: %v", err))

	}
	return &config
}

func (config *Config) RedisConfig() map[string]string {
	return map[string]string{
		"server":   config.Redis.Domain,
		"database": config.Redis.Database,
		"pool":     config.Redis.Pool,
		"process":  config.Redis.Process,
	}
}

func (config *Config) PostgresConfig() string {
	return fmt.Sprintf("sslmode=disable host=%s port=%s dbname=%s user=%s password=%s",
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.Database,
		config.Postgres.Username,
		config.Postgres.Password,
	)
}
