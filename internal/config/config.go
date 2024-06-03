package config

import (
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3"
	"os"
)

type Config struct {
	listen string
	db     string
}

func (c *Config) String() string {
	return fmt.Sprintf("Config(listen=%v dbConnString=%v)", c.listen, "*hide*")
}

func (c *Config) DbConnString() string {
	return c.db
}

func (c *Config) Listen() string {
	return c.listen
}

func Load() (*Config, error) {
	fs := flag.NewFlagSet("cute-chat-backend", flag.ExitOnError)
	cfg := new(Config)
	fs.StringVar(&cfg.listen, "listen", "localhost:46473", "listened http address")
	fs.StringVar(&cfg.db, "db", "", "database connection string")
	fs.String("config", "", "config file (optional)")

	err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("CCB_"),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.JSONParser),
	)
	if err != nil {
		return nil, fmt.Errorf("new config: parse os.Args: %w", err)
	}
	return cfg, nil
}
