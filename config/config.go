package config

import (
	"os"
)

type cfg struct {
	Quoroms    map[string][]string
	Address    string
	ReplFactor int
}

var Config cfg

func GenerateConfig() {
	view := os.Getenv("VIEW")
	addr := os.Getenv("ADDRESS")

	//servers := strings.Split(view, ",")
	Config = cfg{
		//Servers: servers, //TODO update to makeQuorom
		Address: addr,
	}
}

func IsIPInternal(unknownIP string) bool {
	for _, ips := range Config.Quoroms {
		for _, ip := range ips {
			if ip == unknownIP {
				return true
			}
		}
	}
	return false
}
