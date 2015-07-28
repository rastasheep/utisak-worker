package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

type Config struct {
	ReadabilityToken string `json:"readability_token"`
	RedisDomain      string
	Redis            struct {
		Domain   string // location of redis instance
		Database string // instance of the database
		Pool     string // number of connections to keep open with redis
		Process  string // unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
	}
}

var configPath = flag.String("config", "config/config.json", "Path to file containing application ids and credentials for other services.")

func LoadConfig() *Config {
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
