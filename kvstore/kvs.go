package kvstore

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

//Data is the json repsonse
type Data struct {
	Value string `json:"value"`
}

//Message is a message to output when valid
type Message struct {
	Message  string `json:"message"`
	Replaced bool   `json:"replaced"`
}

//ErrorMessage is a message to output when errors occur
type ErrorMessage struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

//GetMessage handles all get messages
type GetMessage struct {
	Exists  bool   `json:"doesExist"`
	Message string `json:"message"`
	Value   string `json:"value"`
}

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
	errMsg := ErrorMessage{}

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

		if err != nil {
			if data.Value == "" {
				errMsg.Error = "Value is missing"
				errMsg.Message = "Error in PUT"
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(errMsg)
			} else {
				valStr, err := s.DAL().Put(key, data.Value)
				msg := Message{}

				if err != nil {
					panic(err)
				} else {
					if valStr == "added" {
						msg.Message = "Added successfully"
						msg.Replaced = false
						w.WriteHeader(http.StatusCreated)
					} else if valStr == "updated" {
						msg.Message = "Updated successfully"
						msg.Replaced = true
						w.WriteHeader(http.StatusOK)
					}
					json.NewEncoder(w).Encode(msg)
				}
			}
		}
	}
}

//GetHandler here
func (s *Store) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	getMsg := GetMessage{}
	getMsg.Exists = ok
	if !ok {
		getMsg.Value = "Key does not exist"
		getMsg.Message = "Error in GET"
		w.WriteHeader(http.StatusNotFound)
	} else {
		val, err := s.DAL().Get(key)
		if err == nil {
			panic(err)
		}
		getMsg.Value = val
		getMsg.Message = "Retrieved successfully"
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(getMsg)
}
