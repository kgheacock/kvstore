package kvstore

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

//DeleteHandler here
func (s *Store) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		w.Write([]byte("something something bad request"))
	}
	s.DAL().Delete(key)
}

//PutHandler puts a key into the store
func (s *Store) PutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	key, ok := vars["key"]
	errMsg := PutFailure{}
	putMsg := PutSuccess{}

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No key"))
	} else if len(key) > 50 {
		errMsg.Error = "Key is too long"
		errMsg.Message = "Error in PUT"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMsg)
	} else {
		data := Data{}
		err := decoder.Decode(&data)

		if err == nil {
			if data.Value == "" {
				errMsg.Error = "Value is missing"
				errMsg.Message = "Error in PUT"
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(errMsg)
			} else {
				valStr, err := s.DAL().Put(key, data.Value)

				if err != nil {
					panic(err)
				} else {
					if valStr == "added" {
						putMsg.Message = "Added successfully"
						putMsg.Replaced = false
						w.WriteHeader(http.StatusCreated)
					} else if valStr == "updated" {
						putMsg.Message = "Updated successfully"
						putMsg.Replaced = true
						w.WriteHeader(http.StatusOK)
					}
					json.NewEncoder(w).Encode(putMsg)
				}
			}
		}
	}
}

//GetHandler here
func (s *Store) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	getMsg := GetSuccess{}
	errMsg := GetFailure{}

	if !ok {
		w.Write([]byte("Method GET not supported"))
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		val, err := s.DAL().Get(key)
		if err != nil {
			errMsg.Exists = false
			errMsg.Error = "Key does not exist"
			errMsg.Message = "Error in GET"
			w.WriteHeader(http.StatusNotFound)
		} else {
			getMsg.Exists = true
			getMsg.Value = val
			getMsg.Message = "Retrieved successfully"
			w.WriteHeader(http.StatusOK)
		}
	}
	json.NewEncoder(w).Encode(getMsg)
}
