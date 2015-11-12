package worker

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
	BaseUrl          string
	ArticlePrefix    string
	ReadabilityToken string
	FeedRegistryPath string
	Postgres         struct {
		Host     string
		Port     string
		Database string
		Username string
		Password string
	}
	Swiftype struct {
		Enabled      bool
		AuthToken    string
		Engine       string
		DocumentType string
	}
}

func LoadConfig() *Config {
	flag.Parse()
	var config Config

	LoadFile(*configPath, &config)
	return &config
}

func LoadFile(path string, dest interface{}) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Could not read config file: %v", err))
	}

	if err = json.Unmarshal(raw, &dest); err != nil {
		panic(fmt.Sprintf("Could not unmarshal config file: %v", err))
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
