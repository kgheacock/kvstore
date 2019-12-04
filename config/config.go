package config

import (
	"os"
	"sync"
)

type cfg struct {
	Quoroms    map[string][]string
	ThisQuorom string
	Address    string
	ReplFactor int
	Mux        sync.Mutex
}

var Config cfg

func GenerateConfig() {
	addr := os.Getenv("ADDRESS")
	Config = cfg{Address: addr, ThisQuorom: "", Quoroms: make(map[string][]string), ReplFactor: 0}
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
