package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/colbyleiske/cse138_assignment2/shard"
	"github.com/colbyleiske/cse138_assignment2/vectorclock"
)

type cfg struct {
	Shards       map[int]shard.Shard
	CurrentShard shard.Shard
	VectorClock  *vectorclock.VectorClock
	Address      string
	ReplFactor   int
	Mux          sync.Mutex
}

var Config cfg

func GenerateConfig() {
	addr := os.Getenv("ADDRESS")
	view := os.Getenv("VIEW")
	replFactor := os.Getenv("REPL_FACTOR")
	replFactorNum, err := strconv.Atoi(replFactor)
	if err != nil {
		log.Fatal("Could not parse replication factor:", replFactor)
	}

	log.Printf("ADDR: %v - VIEW: %v - REPL_FACTOR: %v\n", addr, view, replFactorNum)

	servers := strings.Split(view, ",") // keep all this here lol - we need this to handle making our shards / groups of replicas
	Config = cfg{Address: addr, ReplFactor: replFactorNum}
	Config.Shards = make(map[int]shard.Shard)

	for i := 0; i < len(servers)/replFactorNum; i++ {
		Config.Shards[i] = shard.Shard{ID: i, Nodes: servers[0+(replFactorNum*i) : replFactorNum+(replFactorNum*i)]}
	}
	log.Println(Config.Shards)

}

func IsIPInternal(unknownIP string) bool {
	for _, shards := range Config.Shards {
		for _, ip := range shards.Nodes {
			if ip == unknownIP {
				return true
			}
		}
	}
	return false
}
