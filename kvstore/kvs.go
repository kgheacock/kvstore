package kvstore

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/colbyleiske/cse138_assignment2/config"
	"github.com/gorilla/mux"
)

func (s *Store) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if err := s.DAL().Delete(key); err != nil {
		resp := DeleteResponse{ResponseMessage{"Key does not exist", "Error in DELETE", ""}, false}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := DeleteResponse{ResponseMessage{"", "Deleted successfully", ""}, true}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *Store) KeyCountHandler(w http.ResponseWriter, r *http.Request) {
	//Return Key Count
}

func (s *Store) PutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	decoder := json.NewDecoder(r.Body)

	var data Data
	if err := decoder.Decode(&data); err != nil || data.Value == "" {
		resp := ResponseMessage{Error: "Value is missing", Message: "Error in PUT"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	putResp, err := s.DAL().Put(key, data.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if putResp == ADDED {
		resp := PutResponse{ResponseMessage{"", "Added successfully", ""}, false}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
		return
	}
	if putResp == UPDATED {
		resp := PutResponse{ResponseMessage{"", "Updated successfully", ""}, true}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
}

func (s *Store) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	val, err := s.DAL().Get(key)
	if err != nil {
		resp := GetResponse{ResponseMessage{"Key does not exist", "Error in GET", ""}, false}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := GetResponse{ResponseMessage{"", "Retrieved successfully", val}, true}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
func (s *Store) VCFinishedAckHandler(w http.ResponseWriter, r *http.Request) {
	//ISSUE: This is NOT Itempotent. If a server acks more than once this number will be invalid
	s.nodeAckCount += 1
	if s.nodeAckCount == (s.nodeCount - 1) {
		s.nodeAckCount = 0
		s.viewChangeAllAcksRecievedChannel <- true
	}
}
func (s *Store) ReshardHandler(w http.ResponseWriter, r *http.Request) {
	//source, ok := ctx.Value(middleware.ContextSourceKey).(string)
	vars := mux.Vars(r)
	serverListString := vars["view"]
	serverList := strings.Split(serverListString, ",")
	sort.Strings(serverList)
	for _, server := range serverList {
		config.Config.Servers = append(config.Config.Servers, server)
	}
	/*
		if source == middleware.EXTERNAL {
			s.state = RECIEVED_EXTERNAL_RESHARD
			for _, server := range config.Config.Servers {
				if server != config.Config.addr {
					//TODO
					//sendVCRequest(node, newNodeList)
				}

			}
		}
	*/
	//LOCK server to external requests
	s.state = RECIEVED_INTERNAL_RESHARD
	newServers := s.NewServersFromVC(config.Config.Servers)
	for _, server := range newServers {
		s.hasher.DAL().AddServer(server)

	}
	s.state = TRANSFER_KEYS
	keyList, err := s.dal.KeyList()
	if err != nil {
		log.Fatal(err)
	}
	for _, key := range keyList {
		client := &http.Client{}
		serverIP, err := s.hasher.DAL().ServerOfKey(key)
		if err != nil {
			log.Fatal(err)
		}
		if serverIP != config.Config.Address {
			value, err := s.dal.Get(key)
			if err != nil {
				log.Fatal(err)
			}
			payload, _ := json.Marshal(Data{value})
			req, _ := http.NewRequest("POST", serverIP, bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			//TODO: Add value to JSON payload
			//Make http PUT request to serverIP
			err = s.dal.Delete(key)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	s.state = FINISHED_TRANSFER
	//TODO broadcast ack to all servers
	for _, server := range config.Config.Servers {
		err, _ := http.NewRequest("POST", server, nil)
		if err != nil {
			log.Printf("Ack to server %s returned an error ", server)
		}
	}
	s.ViewChangeFinishedChannel <- true
	s.state = PROCESS_BACKLOG
	<-s.viewChangeAllAcksRecievedChannel
	s.state = NORMAL
}

//Assumes only an INCREASING server list. Invalid when nodes can be deleted
func (s *Store) NewServersFromVC(newNodeList []string) []string {
	oldNodeList := s.hasher.DAL().Servers()
	sort.Strings(oldNodeList)
	var oldNodeID int
	oldNodeID = 0
	var newNodes []string
	for _, node := range newNodeList {
		if node == oldNodeList[oldNodeID] {
			oldNodeID++
		} else {
			newNodes = append(newNodes, node)
		}
	}
	return newNodes
}
