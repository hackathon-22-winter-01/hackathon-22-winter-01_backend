package config

import "flag"

var (
	Port       = flag.Int("port", 8080, "port number")
	Production = flag.Bool("production", false, "enable production mode (default false)")
)

func ParseFlags() {
	flag.Parse()
}
