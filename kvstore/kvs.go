package kvstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/colbyleiske/cse138_assignment2/ctx"
	"github.com/colbyleiske/cse138_assignment2/vectorclock"

	"github.com/colbyleiske/cse138_assignment2/config"
	"github.com/gorilla/mux"
)

func (s *Store) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	var addr string
	if len(r.Header.Get("X-Real-Ip")) != 0 {
		addr = config.Config.Address
	}
	if err := s.DAL().Delete(key); err != nil {
		resp := DeleteResponse{ResponseMessage{"Key does not exist", "Error in DELETE", "", addr, config.Config.CurrentShard().VectorClock}, false}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := DeleteResponse{ResponseMessage{"", "Deleted successfully", "", addr, config.Config.CurrentShard().VectorClock}, true}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
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

	var data Data
	if err := decoder.Decode(&data); err != nil || data.Value == "" {
		resp := ResponseMessage{Error: "Value is missing", Message: "Error in PUT", Address: addr, CausalContext: config.Config.CurrentShard().VectorClock}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	putResp, err := s.DAL().Put(key, data.Value)
	if err != nil {
		log.Printf("Error in PUT. key: %s, value: %s\n", key, data.Value)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	incClock, ok := r.Context().Value(ctx.ContextCausalContextKey).(vectorclock.VectorClock)
	if !ok {
		log.Println("Could not get context from incoming request")
		return
	}
	config.Config.CurrentShard().VectorClock.ReceiveEvent(&incClock)

	if putResp == ADDED {
		resp := PutResponse{ResponseMessage{"", "Added successfully", "", addr, config.Config.CurrentShard().VectorClock}, false}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
		return
	}
	if putResp == UPDATED {
		resp := PutResponse{ResponseMessage{"", "Updated successfully", "", addr, config.Config.CurrentShard().VectorClock}, true}
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

	val, err := s.DAL().Get(key)
	if err != nil {
		resp := GetResponse{ResponseMessage{"Key does not exist", "Error in GET", "", addr, config.Config.CurrentShard().VectorClock}, false}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := GetResponse{ResponseMessage{"", "Retrieved successfully", val, addr, config.Config.CurrentShard().VectorClock}, true}
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
			url := fmt.Sprintf("http://%s/kv-store/keys/%s", serverIP, key)
			value, _ := s.DAL().Get(key)
			data := Data{Value: value}
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
		log.Println()
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
	ctx = context.Background()
	ctx, _ = context.WithTimeout(ctx, config.Config.TimeOut)
	ack = BroadcastMessageAndWait(viewChangeRequest.View, vcBytes, "http://%s/internal/view-change", ctx)
	if !ack {
		log.Println("Recieved 1 or more nacks to view-change request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// //TODO package up the key counts

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
		//var status NodeStatus
		//decoder := json.NewDecoder(resp.Body)
		//err = decoder.Decode(status)
		if err != nil {
			log.Println("Error: 22", err)
		}
		//clusterStatus[status.ShardID] = append(clusterStatus[status.ShardID], status)
	}

}

func (s *Store) GetKeyCountHandler(w http.ResponseWriter, r *http.Request) {
	count := s.DAL().GetKeyCount()
	resp := GetKeyCountRepsponse{"Key count retrieved successfully", count}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
