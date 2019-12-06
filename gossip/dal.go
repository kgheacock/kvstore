package gossip

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/colbyleiske/cse138_assignment2/config"
	"github.com/colbyleiske/cse138_assignment2/vectorclock"
)

//*************** GossipData ***************\\

//GossipData is the data to be sent to another server, key, value, and clock of that request
type GossipData struct {
	VC    map[string]int `json:"vc"`
	Key   string         `json:"key"`
	Value string         `json:"value"`
}

//NewGossipData contains neccessary gossip info
func NewGossipData(key, val string, clock *vectorclock.VectorClock) *GossipData {
	return &GossipData{VC: clock.Clocks, Key: key, Value: val}
}

//*************** GossipQueue ***************\\

//GossipQueue holds queue for requests
type GossipQueue struct {
	Queue []GossipData
}

//NewGossipQueue returns GossipState object
func NewGossipQueue() *GossipQueue {
	q := make([]GossipData, 0, 10)
	return &GossipQueue{Queue: q}
}

//*************** AckTable ***************\\

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

//receivedAllAcks keeps track of ACK's recieved by the servers
func (t *AckTable) receivedAllAcks() bool {
	for _, v := range t.table {
		if v == 0 {
			return false
		}
	}
	return true
}

//*************** Main Gossip Functions ***************\\

//WakeUp starts a ShareGossip request in a bounded time range, forever
func (q *GossipQueue) WakeUp() {
	min := 500
	max := 2000
	ackTable := NewAckTable()

	for {
		rand := rand.Intn(max-min) + min
		time.Sleep(time.Duration(rand) * time.Millisecond)

		if len(q.Queue) > 0 {
			data := (q.Queue)[0]
			ackTable.shareGossip(data)
			//Pop once we received all ACKS
			if ackTable.receivedAllAcks() {
				q.Queue = (q.Queue)[1:]
				//Clear table for garbage collection
				ackTable.table = nil
			}
		}
	}
}

//PrepareForGossip called at PUT endpoint, adds request to GossipQueue
func (q *GossipQueue) PrepareForGossip(key, value string, vc *vectorclock.VectorClock) {
	datagram := NewGossipData(key, value, vc)
	q.Queue = append(q.Queue, *datagram)
	go q.WakeUp()
}

//ShareGossip sends key, val, vectorclock to internal endpoint
func (t *AckTable) shareGossip(datagram GossipData) {
	for _, server := range config.Config.CurrentShard().Nodes {
		payload, _ := json.Marshal(datagram)
		client := &http.Client{Timeout: 500 * time.Millisecond}
		req, _ := http.NewRequest("POST", server, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Real-Ip", config.Config.Address)
		resp, _ := client.Do(req)
		if resp.StatusCode == 200 {
			t.table[server]++
		}
	}

}

//ReceivedGossip handles receiving gossip from another server
func receivedGossip() {
	//Grab received data
	//Check against ours if its newer
	//Update or dont do anything
	//Send an ACK
}
