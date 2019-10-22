package kvstore

import (
	"encoding/json"
	"net/http"

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

func (s *Store) PutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	key := vars["key"]

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
	} else if putResp == UPDATED {
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
