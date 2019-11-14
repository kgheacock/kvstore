package config

import (
	"strings"
)

type cfg struct {
	Servers []string
	Address string
}

var Config cfg

func GenerateConfig() {
	view := os.Getenv("VIEW")
	addr := os.Getenv("ADDRESS")

	servers := strings.Split(view, ",")
	Config = cfg{
		Servers: servers,
		Address: addr,
	}
}
