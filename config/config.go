package config

import (
	"os"
	"strings"
)

type cfg struct {
	Servers    []string
	Address    string
	ReplFactor int
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

func IsIPInternal(unknownIP string) bool {
	for _, ip := range Config.Servers {
		if ip == unknownIP {
			return true
		}
	}
	return false
}
