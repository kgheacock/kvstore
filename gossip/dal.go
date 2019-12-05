package gossip

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/colbyleiske/cse138_assignment2/config"
)

//GossipData is the data to be sent to another server, key, value, and clock of that request
type GossipData struct {
	vc    map[string]int `json:"vc"`
	key   string         `json:"key"`
	value string         `json:"value"`
}

func NewGossipData(k, val, string, clock VectorClock) GossipData {
	return &GossipData{vc: clock, key: k, value: val}
}

//GossipQueue holds queue for requests
type GossipQueue struct {
	Queue []GossipData
}

//NewGossipQueue returns GossipState object
func NewGossipQueue() *GossipQueue {
	//Len of 0, arbitrary capacity of 10
	q := make([]GossipData, 0, 10)
	return &GossipQueue{Queue: q}
}

//WakeUp starts a ShareGossip request in a bounded time range, forever
func (q *GossipQueue) WakeUp() {
	min := 500
	max := 2000
	//To always be listening
	for {
		rand := rand.Intn(max-min) + min
		time.Sleep(time.Duration(rand) * time.Millisecond)
		//To be killed on ACK
		for {
			if len(q.Queue) > 0 {
				x, q := (q.Queue)[0], (q.Queue)[1:]
				ShareGossip(x)
			}

		}
	}
}

//ShareGossip sends key, val, vectorclock to internal endpoint
func ShareGossip(datagram GossipData) {
	gd := GossipData{vc: s.vectorClock().DAL().VC, data: s.DAL().Store}
	payload, _ := json.Marshal(gd)
	client := &http.Clinet{}
	serverIp := config.Config.Address
	req, _ := http.NewRequest("POST", serverIP, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Real-Ip", config.Config.Address)
	resp, err := client.Do(req)
}

//ReceivedGossip handles receiving ONLY gossip from another server
func ReceivedGossip() {
	//Grab received data
	//Check against ours if its newer
	//Update or dont do anything
	//Send an ACK
}

//PrepareForGossip called at endpoint on a Put to server, add to GossipQueue
func (q *GossipQueue) PrepareForGossip(key, value string, vc *VectorClock) {
	//Package up DataEngram
	datagram := NewGossipData(key, value, vc)
	//Add to Queue of requests
	q.Queue = append(q.Queue, datagram)
}
