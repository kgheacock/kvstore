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

//NewGossipData contains neccessary gossip info
func NewGossipData(k, val, string, clock VectorClock) *GossipData {
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

//AckTable holds a map that keeps track of ACK's recieved
type AckTable struct {
	table map[string]int
}

//NewAckTable creates an AckTable object
func NewAckTable() *AckTable {
	numServers := len(config.Config.CurrentShard().Nodes)
	servers := config.Config.CurrentShard().Nodes
	m := make(map[string]int, numServers)
	for _, server := range servers {
		m[server] = 0
	}
	return &AckTable{table: m}
}

//WakeUp starts a ShareGossip request in a bounded time range, forever
func (q *GossipQueue) WakeUp() {
	min := 500
	max := 2000

	for {
		rand := rand.Intn(max-min) + min
		time.Sleep(time.Duration(rand) * time.Millisecond)

		if len(q.Queue) > 0 {
			ackTable := NewAckTable()
			data := (q.Queue)[0]
			ShareGossip(data)
			//Pop once we received all ACKS
			if ackTable.receivedAllAcks() {
				x, q := (q.Queue)[0], (q.Queue)[1:]
			}

		}

	}
}

//receivedAllAcks keeps track of ACK's recieved by the servers
func (t *AckTable) receivedAllAcks() bool {
	for _, v := range t.table {
		if v == 0 {
			return false
		}
	}
	return true
}

//PrepareForGossip called at put endpoint, adds request to GossipQueue
func (q *GossipQueue) PrepareForGossip(key, value string, vc *VectorClock) {
	datagram := NewGossipData(key, value, vc)
	q.Queue = append(q.Queue, *datagram)
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
