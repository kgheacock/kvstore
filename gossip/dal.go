package gossip

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/colbyleiske/cse138_assignment2/config"
)

type GossipData struct {
	vc   []int             `json:"vc"` //Probably make a map instead
	data map[string]string `json:"data"`
}

//WakeUp starts a ShareGossip request in a bounded time range, forever
func WakeUp(s *Store) {
	ServerVC := s.vectorClock.DAL().VC
	min := 500
	max := 2000
	for {
		rand := rand.Intn(max-min) + min
		time.Sleep(time.Duration(rand) * time.Millisecond)
		ShareGossip()
	}
}

//Actually send data
func ShareGossip(s *Store) {
	//INCREMENT VC
	gd := GossipData{vc: s.vectorClock().DAL().VC, data: s.DAL().Store}
	payload, _ := json.Marshal(gd)
	client := &http.Clinet{}
	serverIp := config.Config.Address
	req, _ := http.NewRequest("POST", serverIP, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Real-Ip", config.Config.Address)
	resp, err := client.Do(req)
}

//Handle receiving
func ReceivedGossip() {

}
