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
		Domain   string // location of redis instance
		Database string // instance of the database
		Pool     string // number of connections to keep open with redis
		Process  string // unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
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
