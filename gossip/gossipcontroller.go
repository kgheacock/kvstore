package gossip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"log"
	"net/http"
	"sync"

	"github.com/colbyleiske/cse138_assignment2/config"
)

//GossipData is the data to be sent to another server, key, value, and clock of that request
type GossipItem struct {
	Key        string                        `json:"-"` // gets passed as URL param
	Value      string                        `json:"value"`
	Clock      int                           `json:"lamportclock"`
	Acks       map[string]GossipItemAckState `json:"-"` //map of replica IP to whether its been acked 0 for nothing ,1 for sent, 2 for ack
	AckCount   int                           `json:"-"`
	ID         string                        `json:"-"`
	AckMutex   *sync.Mutex                   `json:"-"`
	CountMutex *sync.Mutex                   `json:"-"`
}

//GossipQueue holds queue for requests
type GossipController struct {
	GossipList   map[uint32]*GossipItem
	IsRunning    bool
	ListMutex    *sync.Mutex
	RunningMutex *sync.Mutex
}

type GossipItemAckState int

const (
	UNACKNOWLEDGED GossipItemAckState = iota
	PENDING
	ACKNOWLEDGED
)

//NewGossipQueue returns GossipState object
func NewGossipController() *GossipController {
	return &GossipController{GossipList: make(map[uint32]*GossipItem), ListMutex: &sync.Mutex{}, IsRunning: false, RunningMutex: &sync.Mutex{}}
}

//WakeUp starts a ShareGossip request in a bounded time range, forever
func (gc *GossipController) StartGossip() {
	log.Println("Starting to gossip...")

	hasGossip := true
	for hasGossip {
		gc.ListMutex.Lock()
		removalList := []uint32{}
		removalMutex := &sync.Mutex{}
		var wg sync.WaitGroup

		for hashedKey, item := range gc.GossipList {
			wg.Add(1)
			go func(hashedKey uint32, item *GossipItem, wg *sync.WaitGroup) {
				item.AckMutex.Lock()
				for node, acked := range item.Acks {
					if acked == UNACKNOWLEDGED {
						item.Acks[node] = PENDING
						go item.AttemptRequest(node)
					}
				}
				item.AckMutex.Unlock()
				item.CountMutex.Lock()
				if item.AckCount >= config.Config.ReplFactor {
					removalMutex.Lock()
					removalList = append(removalList, hashedKey)
					removalMutex.Unlock()
				}
				item.CountMutex.Unlock()
				wg.Done()
			}(hashedKey, item, &wg)
		}
		gc.ListMutex.Unlock()
		wg.Wait()

		removalMutex.Lock()
		gc.ListMutex.Lock()
		for _, i := range removalList {
			delete(gc.GossipList, i)
		}
		removalMutex.Unlock()

		hasGossip = !(len(gc.GossipList) == 0)
		gc.ListMutex.Unlock()
	}

	log.Println("Finished gossiping...")
	gc.RunningMutex.Lock()
	gc.IsRunning = false
	gc.RunningMutex.Unlock()
}

//PrepareForGossip called at PUT endpoint, adds request to GossipQueue
func (gc *GossipController) AddGossipItem(key string, val string, clock int) {
	log.Println("Added item to gossip")
	acks := make(map[string]GossipItemAckState)
	for _, ip := range config.Config.CurrentShard().Nodes {
		acks[ip] = UNACKNOWLEDGED
	}
	acks[config.Config.Address] = ACKNOWLEDGED

	//AckCount is set to 1 to simulate ourselves acking a request
	item := &GossipItem{Clock: clock, Key: key, Value: val, AckCount: 1, Acks: acks, AckMutex: &sync.Mutex{}, CountMutex: &sync.Mutex{}}

	gc.ListMutex.Lock()
	hashedKey := crc32.ChecksumIEEE([]byte(fmt.Sprintf("%v%v%v%v", item.AckCount, item.Clock, item.Key, item.Value)))
	gc.GossipList[hashedKey] = item
	gc.ListMutex.Unlock()

	gc.RunningMutex.Lock()
	if !gc.IsRunning {
		gc.IsRunning = true
		go gc.StartGossip()
	}
	gc.RunningMutex.Unlock()
}

func (i *GossipItem) AttemptRequest(node string) {
	marshalledItem, err := json.Marshal(i)
	if err != nil {
		log.Println("error marshalling gossip item", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Config.TimeOut)
	defer cancel()
	client := &http.Client{}
	addr := fmt.Sprintf("http://%s/internal/gossip-put/%s", node, i.Key)
	req, err := http.NewRequest("PUT", addr, bytes.NewBuffer(marshalledItem))
	if err != nil {
		log.Println("error creating request", err)
		i.AckMutex.Lock()
		i.Acks[node] = UNACKNOWLEDGED
		i.AckMutex.Unlock()
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Real-Ip", config.Config.Address)
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		log.Println(err)
		i.AckMutex.Lock()
		i.Acks[node] = UNACKNOWLEDGED
		i.AckMutex.Unlock()
		return //we can assume they did not update their kvs
	}

	if resp.StatusCode == http.StatusOK {
		//No idea proper use of mutexes and this is really ugly looking but i wanted to be explicit for a first time writing this stuff
		i.AckMutex.Lock()
		i.Acks[node] = ACKNOWLEDGED
		i.AckMutex.Unlock()
		i.CountMutex.Lock()
		i.AckCount++
		i.CountMutex.Unlock()
	}
}

func (gc *GossipController) IsGossiping() (running bool) {
	gc.RunningMutex.Lock()
	running = gc.IsRunning
	gc.RunningMutex.Unlock()
	return running
}
