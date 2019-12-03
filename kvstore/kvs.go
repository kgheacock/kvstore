package kvstore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	respStruct := shard{KeyCount: s.DAL().GetKeyCount(), Address: config.Config.Address}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respStruct)
	s.state = NORMAL
}
func (s *Store) InternalReshardHandler(w http.ResponseWriter, r *http.Request) {
	s.state = RECIEVED_INTERNAL_RESHARD
	var viewChangeRequest InternalViewChangeRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&viewChangeRequest); err != nil || viewChangeRequest.ReplFactor == 0 {
		log.Println()
		log.Println("ERROR: 1", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for quorom := range config.Config.Quoroms {
		s.hasher.DAL().RemoveServer(quorom)
	}
	config.Config.Mux.Lock()
	config.Config.Quoroms = viewChangeRequest.NamedView
	config.Config.Mux.Unlock()
	for quorom, servers := range config.Config.Quoroms {
		s.hasher.DAL().RemoveServer(quorom)
		for _, server := range servers {
			if server == config.Config.Address {
				config.Config.ThisQuorom = quorom
			}
		}
	}
	client := &http.Client{}
	keyList := s.DAL().KeyList()
	for _, key := range keyList {
		properQuorom, err := s.hasher.DAL().GetServerByKey(key)
		if err != nil {
			log.Println("ERROR: 13", err)
		}
		if properQuorom != config.Config.ThisQuorom {
			url := fmt.Sprintf("http://%s/kv-store/keys/%s", config.Config.Quoroms[properQuorom][0], key)
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
func (s *Store) ExternalReshardHandler(w http.ResponseWriter, r *http.Request) {
	//First to recieve. We must forward to everyone but ourselves. Then call reshard ourselves
	var viewChangeRequest ExternalViewChangeRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&viewChangeRequest); err != nil || viewChangeRequest.ReplFactor == 0 {
		log.Println()
		log.Println("ERROR: 1", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	viewChangeSentChannel := make(chan bool, len(viewChangeRequest.View))
	replFacotr := viewChangeRequest.ReplFactor
	var namedQuorum map[string][]string = nil //:=makeQuorum(config.Config.Servers, replFactor)
	namedViewChangeRequest := InternalViewChangeRequest{NamedView: namedQuorum, ReplFactor: replFacotr}
	viewChangeReqBytes, err := json.Marshal(namedViewChangeRequest)
	if err != nil {
		log.Println("ERROR: 10", err)
	}
	for _, server := range viewChangeRequest.View {
		go func() {
			client := &http.Client{}
			url := fmt.Sprintf("http://%s/internal/view-change", server)
			req, err := http.NewRequest("PUT", url, bytes.NewReader(viewChangeReqBytes))
			if err != nil {
				log.Println("ERROR: 2", err)
				w.WriteHeader(http.StatusInternalServerError)
				viewChangeSentChannel <- false
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Real-Ip", config.Config.Address)
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				viewChangeSentChannel <- false
				return
			}
			resp.Body.Close()
			viewChangeSentChannel <- true
		}()
	}
	for _, server := range viewChangeRequest.View {
		vcAck := <-viewChangeSentChannel
		if !vcAck {
			log.Printf("Recieved an error response from server %s", server)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	//TODO package up the key counts
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
	}

}

func (s *Store) GetKeyCountHandler(w http.ResponseWriter, r *http.Request) {
	count := s.DAL().GetKeyCount()
	resp := GetKeyCountRepsponse{"Key count retrieved successfully", count}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
