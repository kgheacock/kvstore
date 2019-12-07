package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/colbyleiske/cse138_assignment2/shard"
	"github.com/colbyleiske/cse138_assignment2/vectorclock"
)

type cfg struct {
	Shards         map[int]*shard.Shard
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
	replFactorNum, err := strconv.Atoi(replFactor)
	if err != nil {
		log.Fatal("Could not parse replication factor:", replFactor)
	}

	log.Printf("ADDR: %v - VIEW: %v - REPL_FACTOR: %v\n", addr, view, replFactorNum)

	//Ideally I move all this shard bullshit over to the shard package later - someone remind me pls <3

	servers := strings.Split(view, ",") // keep all this here lol - we need this to handle making our shards / groups of replicas
	Config = cfg{Address: addr, ReplFactor: replFactorNum, TimeOut: 5 * time.Second}
	Config.Shards = make(map[int]*shard.Shard)

	for i := 0; i < len(servers)/replFactorNum; i++ {
		Config.Shards[i] = &shard.Shard{ID: strconv.Itoa(i), Nodes: servers[0+(replFactorNum*i) : replFactorNum+(replFactorNum*i)]}
		if Contains(servers[0+(replFactorNum*i):replFactorNum+(replFactorNum*i)], addr) {
			//Only need to set the vector clock on OUR shard - don't care about the other shards
			Config.Shards[i].VectorClock = vectorclock.NewVectorClock(Config.Shards[i].Nodes, addr)
			Config.CurrentShardID = i
		}
	}

	log.Println(Config.Shards)
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
