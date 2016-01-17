package worker

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var configPath = flag.String("config", "config/config.json", "Path to configuration file")

type Config struct {
	LogTo         string
	LogLevel      string
	BaseUrl       string
	ArticlePrefix string
	Swiftype      Swiftype
	Feeds         []*Feed
}

type Swiftype struct {
	Enabled      bool
	AuthToken    string
	Engine       string
	DocumentType string
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
		os.Getenv("POSTGRES_SERVICE_HOST"),
		os.Getenv("POSTGRES_SERVICE_PORT"),
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_USERNAME"),
		os.Getenv("POSTGRES_PASSWORD"),
	)
}

func (config *Config) SwiftypeConfig() Swiftype {
	swiftype := config.Swiftype
	swiftype.AuthToken = os.Getenv("SWIFTYPE_TOKEN")
	return swiftype
}
