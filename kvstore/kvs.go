package kvstore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/colbyleiske/cse138_assignment2/config"
	"github.com/colbyleiske/cse138_assignment2/ctx"
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
		resp := DeleteResponse{ResponseMessage{"Key does not exist", "Error in DELETE", "", addr}, false}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := DeleteResponse{ResponseMessage{"", "Deleted successfully", "", addr}, true}
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
		resp := ResponseMessage{Error: "Value is missing", Message: "Error in PUT", Address: addr}
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

	if putResp == ADDED {
		resp := PutResponse{ResponseMessage{"", "Added successfully", "", addr}, false}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
		return
	}
	if putResp == UPDATED {
		resp := PutResponse{ResponseMessage{"", "Updated successfully", "", addr}, true}
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
		resp := GetResponse{ResponseMessage{"Key does not exist", "Error in GET", "", addr}, false}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := GetResponse{ResponseMessage{"", "Retrieved successfully", val, addr}, true}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
func (s *Store) ReshardCompleteHandler(w http.ResponseWriter, r *http.Request) {
	s.state = NORMAL
	respStruct := shard{KeyCount: s.DAL().GetKeyCount(), Address: config.Config.Address}
	//s.ViewChangeFinishedChannel <- true
	//s.ViewChangeFinishedChannel = make(chan bool, 1)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respStruct)
}
func (s *Store) ReshardHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var viewChangeRequest ViewChangeRequest
	if err := decoder.Decode(&viewChangeRequest); err != nil || viewChangeRequest.View == "" {
		return
	}

	serverList := strings.Split(viewChangeRequest.View, ",")
	newServers := s.Difference(serverList, config.Config.Servers)
	config.Config.Servers = serverList

	source, ok := r.Context().Value(ctx.ContextSourceKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to find source of request")
		return
	}
	if source == ctx.EXTERNAL {
		//LOCK server to external requests
		s.state = RECIEVED_EXTERNAL_RESHARD
		//TODO: calls are currently 1 at a time but this can be optimized by
		//making the requests in a GO routine and using a channel to count acks
		//e.g.:
		/*
			ack_count := 0
			ack_channel := make(chan bool, node_count)
			go func() {
				client := &http.Client{}

				resp,err :=//make requet
				for err!=nil{
					resp,err = //resend request but for debugging just do a log.Fatal()
				}
				ack_channel <- true
			 }()

			 for ack_count < (node_count - 1){
				 if <- ack_channel {
					 ack_count ++
				 }
			 }
			 //send VC complete message
		*/
		client := &http.Client{}
		for _, server := range config.Config.Servers {
			log.Println(server)
			if server != config.Config.Address {
				url := fmt.Sprintf("http://%s/kv-store/view-change", server)
				viewChangeReqBytes, err := json.Marshal(viewChangeRequest)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println(err)
					return
				}
				req, err := http.NewRequest("PUT", url, bytes.NewReader(viewChangeReqBytes))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Printf("Error in PUT request. URL: %s", url)
					return
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Real-Ip", config.Config.Address)
				resp, err := client.Do(req)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println(err)
					return
				}
				resp.Body.Close()
			}
		}
	}
	//LOCK server to external requests
	s.state = RECIEVED_INTERNAL_RESHARD
	for _, server := range newServers {
		s.hasher.DAL().AddServer(server)
	}

	s.state = TRANSFER_KEYS
	keyList := s.dal.KeyList()

	for _, key := range keyList {
		client := &http.Client{}
		serverIP, err := s.hasher.DAL().GetServerByKey(key)
		if err != nil {
			log.Printf("Error Getting Server by Key. key: %s\n", key)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if serverIP != config.Config.Address {
			value, err := s.dal.Get(key)
			if err != nil {
				log.Printf("Error on GET. key: %s\n", key)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			url := fmt.Sprintf("http://%s/kv-store/keys/%s", serverIP, key)
			payload, _ := json.Marshal(Data{value})
			req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Real-Ip", config.Config.Address)
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error on Request: %s", req)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()
			log.Printf("Transfered key: %s to server: %s", key, serverIP)
			err = s.dal.Delete(key)
			if err != nil {
				log.Printf("Error on Delete. key: %s\n", key)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	s.state = WAITING_FOR_ACK
	if source == ctx.EXTERNAL {
		client := &http.Client{}
		var shardList []shard
		for _, server := range config.Config.Servers {
			if server != config.Config.Address {
				url := fmt.Sprintf("http://%s/internal/vc-complete", server)
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Real-Ip", config.Config.Address)
				var s shard
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Error on Client Do request. Request: %s\n", req)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				decoder := json.NewDecoder(resp.Body)
				err = decoder.Decode(&s)
				shardList = append(shardList, s)
				if err != nil {
					log.Printf("Error on adding shard. Shard: %s\n", s)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}
		shardList = append(shardList, shard{KeyCount: s.DAL().GetKeyCount(), Address: config.Config.Address})
		resp := struct {
			Message string  `json:"message"`
			Shards  []shard `json:"shards"`
		}{"View change successful", shardList}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		s.state = NORMAL
		//s.ViewChangeFinishedChannel <- true
		//s.ViewChangeFinishedChannel = make(chan bool, 1)
	}
}

func (s *Store) Difference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func (s *Store) GetKeyCountHandler(w http.ResponseWriter, r *http.Request) {
	count := s.DAL().GetKeyCount()
	resp := GetKeyCountRepsponse{"Key count retrieved successfully", count}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
