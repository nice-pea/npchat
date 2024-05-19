package config

import (
	"flag"
	"github.com/peterbourgon/ff/v3"
	"log"
	"os"
)

type Config struct {
	listen string
}

func (c *Config) Listen() string {
	return c.listen
}

func Load() *Config {
	fs := flag.NewFlagSet("cute-chat-backend", flag.ExitOnError)
	cfg := new(Config)
	fs.StringVar(&cfg.listen, "listen", "localhost:46473", "listen listen")
	//refresh    = fs.Duration("refresh", 15*time.Second, "refresh interval")
	//debug      = fs.Bool("debug", false, "log debug information")
	fs.String("config", "", "config file (optional)")

	err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("CCB_"),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
	)
	if err != nil {
		log.Fatalf("new config: parse os.Args: %v", err)
	}
	return cfg
}
