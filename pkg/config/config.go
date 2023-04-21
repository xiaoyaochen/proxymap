package config

import (
	"flag"
)

type Config struct {
	Addr        string // proxy listen addr
	Upstream    string // upstream listen addr
	SslInsecure bool   //Upstream proxy SslInsecure
	DB          string //db save response
}

func LoadConfigFromCli() *Config {
	config := new(Config)
	flag.StringVar(&config.Addr, "addr", ":9080", "proxy listen addr")
	flag.StringVar(&config.Upstream, "upstream", "", "Upstream proxy listen addr")
	flag.BoolVar(&config.SslInsecure, "upssl", true, "Upstream proxy is in InsecureSkipVerify")
	flag.StringVar(&config.DB, "db", "", "db save response")
	flag.Parse()

	return config
}
