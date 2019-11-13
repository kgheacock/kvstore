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
	//view := os.Getenv("VIEW")
	//addr := os.Getenv("ADDRESS")
	addr := "localhost:13800"
	view := "localhost,10.10.0.2:13800" //clTODO
	servers := strings.Split(view, ",")
	Config = cfg{
		Servers: servers,
		Address: addr,
	}
}
