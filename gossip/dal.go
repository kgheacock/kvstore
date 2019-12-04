package gossip

import (
	"time"
	"math/rand"
)

type GossipData struct{
	vc []int `json:"vc"` //Probably make a map instead
	data map[string]string `json:"data"`
}

//WakeUp starts a ShareGossip request in a bounded time range, forever
func WakeUp() {
	ServerVC := NewVectorClock()
	min := 500
	max := 2000
	for {
		rand := rand.Intn(max - min) + min
		time.Sleep(time.Duration(rand) * time.Millisecond)
		ShareGossip()
	}
}

//Actually send data
func ShareGossip() {
	//INCREMENT VC
	gd := GossipData{vc:...,data:...}
	payload,_ := json.Marshal(gd)
	client := &http.Clinet{}
	serverIp := \\ip of server
	req, _ := http.NewRequest("POST", serverIP, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Real-Ip",config.Config.Address)
	resp, err := client.Do(req)
}

//Handle receiving
func ReceivedGossip() {

}
