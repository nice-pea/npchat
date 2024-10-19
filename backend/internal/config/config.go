package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/peterbourgon/ff/v3"
)

type Config struct {
	Listen     string
	DB         string
	ConfigFile string
}

func parseFile(file string) (cfg Config, err error) {
	var b []byte
	if b, err = os.ReadFile(file); err != nil {
		return Config{}, err
	}

	return cfg, json.Unmarshal(b, &cfg)
}

func Load() (cfg Config, err error) {
	fs := flag.NewFlagSet("nice-pea-chat", flag.ExitOnError)
	fs.StringVar(&cfg.Listen, "Listen", "localhost:46473", "listened http address")
	fs.StringVar(&cfg.DB, "DB", "", "database connection string")
	fs.StringVar(&cfg.ConfigFile, "config", "config.json", "config file (optional)")

	return cfg, ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.JSONParser),
	)
}
