package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/colbyleiske/cse138_assignment2/shard"
)

type cfg struct {
	Shards         map[int]*shard.Shard
	RawServers []string
	CurrentShardID int
	Address        string
	ReplFactor     int
	TimeOut        time.Duration
	Mux            sync.Mutex
}

var Config cfg

func GenerateConfig() {
	addr := os.Getenv("ADDRESS")
	view := os.Getenv("VIEW")
	replFactor := os.Getenv("REPL_FACTOR")
	if replFactor == "" {
		replFactor = "1"
	}
	replFactorNum, err := strconv.Atoi(replFactor)
	if err != nil {
		log.Fatal("Could not parse replication factor:", replFactor)
	}

	log.Printf("ADDR: %v - VIEW: %v - REPL_FACTOR: %v\n", addr, view, replFactorNum)

	//Ideally I move all this shard bullshit over to the shard package later - someone remind me pls <3

	Config = cfg{Address: addr, ReplFactor: replFactorNum, TimeOut: 5 * time.Second}
	Config.RawServers = strings.Split(view, ",") // keep all this here lol - we need this to handle making our shards / groups of replicas
	Config.Shards = make(map[int]*shard.Shard)

	for i := 0; i < len(Config.RawServers)/replFactorNum; i++ {
		Config.Shards[i] = &shard.Shard{ID: strconv.Itoa(i), Nodes: Config.RawServers[0+(replFactorNum*i) : replFactorNum+(replFactorNum*i)]}
		if Contains(Config.RawServers[0+(replFactorNum*i):replFactorNum+(replFactorNum*i)], addr) {
			Config.CurrentShardID = i
		}
	}
}

func (config cfg) CurrentShard() *shard.Shard {
	return config.Shards[config.CurrentShardID]
}

func Contains(servers []string, ip string) bool {
	for _, serverIP := range servers {
		if serverIP == ip {
			return true
		}
	}
	return false
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
