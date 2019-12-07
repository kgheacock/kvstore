package kvstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/colbyleiske/cse138_assignment2/ctx"

	"github.com/colbyleiske/cse138_assignment2/config"
	"github.com/gorilla/mux"
)

func (s *Store) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// key := vars["key"]
	// var addr string
	// if len(r.Header.Get("X-Real-Ip")) != 0 {
	// 	addr = config.Config.Address
	// }
	// if err := s.DAL().Delete(key); err != nil {
	// 	resp := DeleteResponse{ResponseMessage{"Key does not exist", "Error in DELETE", "", addr, config.Config.CurrentShard().VectorClock}, false}
	// 	w.WriteHeader(http.StatusNotFound)
	// 	json.NewEncoder(w).Encode(resp)
	// 	return
	// }

	// resp := DeleteResponse{ResponseMessage{"", "Deleted successfully", "", addr, config.Config.CurrentShard().VectorClock}, true}
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(resp)
}

func (s *Store) KeyCountHandler(w http.ResponseWriter, r *http.Request) {
	//Return Key Count
	resp := struct {
		message  string `json:"message"`
		keyCount int    `json:"key-count"`
	}{message: "Key count retrieved successfully", keyCount: s.DAL().GetKeyCount()}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *Store) PutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	decoder := json.NewDecoder(r.Body)
	var addr string
	if len(r.Header.Get("X-Real-Ip")) != 0 {
		addr = config.Config.Address
	}

	incClock, ok := r.Context().Value(ctx.ContextCausalContextKey).(map[string]int)
	if !ok {
		log.Println("test")
		log.Println("Could not get context from incoming request")
		return
	}

	var data Data
	if err := decoder.Decode(&data); err != nil || data.Value == "" {
		resp := ResponseMessage{Error: "Value is missing", Message: "Error in PUT", Address: addr, CausalContext: incClock}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	curClock, ok := s.DAL().MapKeyToClock()[key]
	if !ok {
		curClock = 0
	}
	if incClock[key] > curClock {
		curClock = incClock[key]
	}

	putResp := s.DAL().Put(key, StoredValue{data.Value, curClock + 1})
	incClock[key] = curClock + 1

	s.gossipController.AddGossipItem(key, data.Value, incClock[key])

	if putResp == ADDED {
		resp := PutResponse{ResponseMessage{"", "Added successfully", "", addr, incClock}, false}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
		return
	}
	if putResp == UPDATED {
		resp := PutResponse{ResponseMessage{"", "Updated successfully", "", addr, incClock}, true}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
}

func (s *Store) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	var addr string
	if len(r.Header.Get("X-Real-Ip")) != 0 {
		addr = config.Config.Address
	}

	incClock, ok := r.Context().Value(ctx.ContextCausalContextKey).(map[string]int)
	if !ok {
		log.Println("Could not get context from incoming request")
		return
	}

	val, err := s.DAL().Get(key)
	if err != nil {
		if incClock[key] == 0 {
			delete(incClock, key)
		}
		resp := GetResponse{ResponseMessage{"Key does not exist", "Error in GET", "", addr, incClock}, false}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	incClock[key] = val.lamportclock

	resp := GetResponse{ResponseMessage{"", "Retrieved successfully", val.value, addr, incClock}, true}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *Store) PrepareReshardHandler(w http.ResponseWriter, r *http.Request) {
	s.state = PREPARE_FOR_RESHARD
	//Gossip to all nodes
	w.WriteHeader(http.StatusOK)
}

func (s *Store) ReshardCompleteHandler(w http.ResponseWriter, r *http.Request) {
	respStruct := NodeStatus{KeyCount: s.DAL().GetKeyCount(), IP: config.Config.Address, ShardID: config.Config.CurrentShardID}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respStruct)
	s.state = NORMAL
}

func (s *Store) InternalReshardHandler(w http.ResponseWriter, r *http.Request) {
	s.state = RECIEVED_INTERNAL_RESHARD
	var viewChangeRequest ViewChangeRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&viewChangeRequest); err != nil || viewChangeRequest.ReplFactor == 0 {
		log.Println()
		log.Println("ERROR: 1", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for shardID := range config.Config.Shards {
		s.hasher.DAL().RemoveShard(config.Config.Shards[shardID])
	}
	config.Config.Mux.Lock()
	config.Config.Shards = makeShards(viewChangeRequest.View, viewChangeRequest.ReplFactor)
	config.Config.ReplFactor = viewChangeRequest.ReplFactor
	for shardID, shard := range config.Config.Shards {
		s.hasher.DAL().AddShard(shard)
		for _, server := range shard.Nodes {
			if server == config.Config.Address {
				config.Config.CurrentShardID = shardID
			}
		}
	}
	config.Config.Mux.Unlock()
	client := &http.Client{}
	keyList := s.DAL().KeyList()
	for _, key := range keyList {
		serverIP, err := s.hasher.DAL().GetServerByKey(key)
		if err != nil {
			log.Println("ERROR: 13", err)
		}
		if serverIP != config.Config.Address {
			url := fmt.Sprintf("http://%s/internal/reshard-put/%s", serverIP, key)
			value, _ := s.DAL().Get(key)
			data := Data{Value: value.value}
			payload, err := json.Marshal(data)
			if err != nil {
				log.Println("ERROR: 5", err)
			}
			req, err := http.NewRequest("PUT", url, bytes.NewReader(payload))
			if err != nil {
				log.Println("ERROR: 6", err)
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Println("ERROR: 7", err, " ", resp)
			}
			resp.Body.Close()
		}
	}

}
func BroadcastMessageAndWait(serverList []string, data []byte, urlFmtString string, ctx context.Context) bool {
	acksRecieved := make(chan bool, len(serverList))
	for _, server := range serverList {
		go func() {
			client := &http.Client{}
			url := fmt.Sprintf(urlFmtString, server)
			req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
			if err != nil {
				log.Println("ERROR: 2", err)
				acksRecieved <- false
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Real-Ip", config.Config.Address)
			req = req.WithContext(ctx)
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				acksRecieved <- false
				return
			}
			resp.Body.Close()
			acksRecieved <- true
		}()
	}
	retVal := true
	for _, server := range serverList {
		ack := <-acksRecieved
		if !ack {
			log.Printf("Recieved an error response from server %s", server)
			retVal = false
		}
	}
	return retVal
}
func (s *Store) ExternalReshardHandler(w http.ResponseWriter, r *http.Request) {
	var viewChangeRequest ViewChangeRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&viewChangeRequest); err != nil || viewChangeRequest.ReplFactor == 0 {
		log.Println("ERROR: 1", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//First to recieve. We must first prepare everyone
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, config.Config.TimeOut)
	ack := BroadcastMessageAndWait(viewChangeRequest.View, nil, "http://%s/internal/prepare-for-vc", ctx)
	if !ack {
		log.Println("Recieved 1 or more nacks to prepare-view-change")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vcBytes, err := json.Marshal(viewChangeRequest)
	if err != nil {
		log.Printf("ERROR: 33", err)
	}
	//After everyone is prepare then we send out the request
	ctx = context.Background()
	ctx, _ = context.WithTimeout(ctx, config.Config.TimeOut)
	ack = BroadcastMessageAndWait(viewChangeRequest.View, vcBytes, "http://%s/internal/view-change", ctx)
	if !ack {
		log.Println("Recieved 1 or more nacks to view-change request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Finally we pack up all responses and send as a response
	type shardStatus struct {
		shardID  int
		keyCount int
		replicas []string
	}
	type vcResponse struct {
		shards        []shardStatus  `json:"shards"`
		message       string         `json:"message"`
		CausalContext map[string]int `json:"causal-context"`
	}
	shardMap := make(map[int]shardStatus)
	for _, server := range viewChangeRequest.View {
		client := &http.Client{}
		url := fmt.Sprintf("http://%s/internal/vc-complete", server)
		req, err := http.NewRequest("GET", url, bytes.NewReader(nil))
		if err != nil {
			log.Println("ERROR: 2", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Real-Ip", config.Config.Address)
		resp, err := client.Do(req)
		if err != nil {
			log.Println("ERROR: 9", err, " ", resp)
		}

		if err != nil {
			log.Println("Error: 22", err)
		}
		var ns NodeStatus
		decode := json.NewDecoder(resp.Body)
		decode.Decode(ns)
		currentShardStatus := shardMap[ns.ShardID]
		currentShardStatus.replicas = append(currentShardStatus.replicas, ns.IP)
		currentShardStatus.keyCount = ns.KeyCount
		currentShardStatus.shardID = ns.ShardID
		shardMap[ns.ShardID] = currentShardStatus
	}
	vcResp := vcResponse{message: "View change successful", CausalContext: make(map[string]int), shards: []shardStatus{}}
	for _, value := range shardMap {
		vcResp.shards = append(vcResp.shards, value)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vcResp)
	return

}

func (s *Store) GetKeyCountHandler(w http.ResponseWriter, r *http.Request) {
	count := s.DAL().GetKeyCount()
	cc, ok := r.Context().Value(ctx.ContextCausalContextKey).(map[string]int)
	if !ok {
		errResp := ResponseMessage{Message: "Error in GET", Error: "Can't receive Causal Context", CausalContext: cc}
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(errResp)
		return
	}

	for key, clock := range cc {
		ourKey := true
		proposedClock, err := s.DAL().Get(key)
		if err != nil { // key does not exist
			if server, _ := s.hasher.DAL().GetServerByKey(key); server == config.Config.Address {
				//our key - we don't have it so we are behind
				proposedClock.lamportclock = 0
			} else {
				ourKey = false // not on our shard - don't bother checking
			}
		}

		if ourKey && proposedClock.lamportclock < clock {
			errResp := ResponseMessage{Message: "Error in GET", Error: "Unable to satisfy request", CausalContext: cc}
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(errResp)
			return
		}
	}

	resp := GetKeyCountRepsponse{"Key count retrieved successfully", count, config.Config.CurrentShardID, cc}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *Store) ReshardPutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	decoder := json.NewDecoder(r.Body)

	data := struct {
		Value        string `json:"value"`
		LamportClock int    `json:"lamportclock"`
	}{}

	if err := decoder.Decode(&data); err != nil || data.Value == "" {
		//Don't need real formatted response since its all internal
		resp := ResponseMessage{Error: "Value is missing", Message: "Error in PUT"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	s.DAL().Put(key, StoredValue{data.Value, data.LamportClock})

	resp := ResponseMessage{Message: "Updated successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return
}

func (s *Store) GossipPutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	decoder := json.NewDecoder(r.Body)

	data := struct {
		Value        string `json:"value"`
		LamportClock int    `json:"lamportclock"`
	}{}

	if err := decoder.Decode(&data); err != nil || data.Value == "" {
		//Don't need real formatted response since its all internal
		resp := ResponseMessage{Error: "Value is missing", Message: "Error in PUT"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	log.Println("received gossip put for", key, data.Value, data.LamportClock)
	proposedKeyClock, ok := s.DAL().MapKeyToClock()[key]
	if !ok {
		proposedKeyClock = 0 //just incase we do not have this value
	}

	if data.LamportClock > proposedKeyClock {
		s.DAL().Put(key, StoredValue{data.Value, data.LamportClock})
	}

	if data.LamportClock == proposedKeyClock {
		newVal := data.Value
		if val, _ := s.DAL().Get(key); val.Value() < newVal { //just arbitrarily pick lower val
			newVal = val.Value()
		}
		s.DAL().Put(key, StoredValue{newVal, data.LamportClock})
	}

	//if the incoming clock for our key is LOWER , we send a 200 and do not operate. We have more context in this case.

	resp := ResponseMessage{Message: "Updated successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return
}

func (s *Store) GetShardByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	checkID := vars["id"]
	shardID, err := strconv.Atoi(checkID)

	cc := make(map[string]int)

	if err != nil {
		errResp := ResponseMessage{
			Message:       "Error in GET",
			Error:         "Can't parse Id",
			CausalContext: cc,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResp)
		return
	}

	replicas, ok := config.Config.Shards[shardID]

	if !ok {
		errResp := ResponseMessage{
			Message:       "Error in GET",
			Error:         "ID not found",
			CausalContext: cc,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResp)
		return
	}

	count := s.DAL().GetKeyCount()

	resp := GetShardByIdResponse{
		ResponseMessage{
			Message:       "Shard information retrieved successfully",
			CausalContext: cc,
		},
		checkID,
		count,
		replicas.Nodes,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *Store) GetShardHandler(w http.ResponseWriter, r *http.Request) {
	shards := make([]int, 0)
	cc, ok := r.Context().Value(ctx.ContextCausalContextKey).(map[string]int)

	if !ok {
		fmt.Println("No clocks")
		errResp := ResponseMessage{Message: "Error in GET", Error: "Can't recieve Causal Context", CausalContext: cc}
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(errResp)
	}

	for key := range config.Config.Shards {
		shards = append(shards, key)
	}

	resp := GetShardResponse{
		ResponseMessage{
			Message:       "Shard membership retrieved successfully",
			CausalContext: cc,
		},
		shards,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}
