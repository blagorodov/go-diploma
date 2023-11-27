package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress  string
	DatabaseURI    string
	AccrualAddress string
	LogPath        string
}

var Options = loadConfig()

func loadConfig() Config {
	o := Config{}
	flag.StringVar(&o.ServerAddress, "a", ":8090", "address and port to run server")
	flag.StringVar(&o.DatabaseURI, "d", "", "database address uri")
	flag.StringVar(&o.AccrualAddress, "r", ":8080", "accrual system address and port")
	flag.StringVar(&o.LogPath, "l", "log.log", "path to log file")
	flag.Parse()

	if e := os.Getenv("RUN_ADDRESS"); e != "" {
		o.ServerAddress = e
	}
	if e := os.Getenv("DATABASE_URI"); e != "" {
		o.DatabaseURI = e
	}
	if e := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); e != "" {
		o.AccrualAddress = e
	}
	return o
}
