package config

import (
	"encoding/json"
	"flag"
	"os"
)

type Config struct {
	App struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"app"`
	Database struct {
		Url string `json:"url"`
	} `json:"database"`
}

func Load() (Config, error) {
	var file string
	fs := flag.NewFlagSet("nice-pea-chat", flag.ExitOnError)
	fs.StringVar(&file, "config", "config.json", "config file in json format")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return Config{}, err
	}

	return parseFile(file)
}

func parseFile(file string) (cfg Config, err error) {
	var b []byte
	if b, err = os.ReadFile(file); err != nil {
		return Config{}, err
	}

	return cfg, json.Unmarshal(b, &cfg)
}
