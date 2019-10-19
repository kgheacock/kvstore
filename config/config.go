package config

import (
	"os"
)

type cfg struct {
	IsFollower bool
	ForwardAddress string
}

var Config cfg 

func GenerateConfig() {
	addr := os.Getenv("FORWARDING_ADDRESS")
	Config = cfg {
		IsFollower: len(addr) > 0,
		ForwardAddress: addr,
	}
}