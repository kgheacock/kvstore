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

//PutHandler here
func (s *Store) PutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	key, ok := vars["key"]
	if !ok {
		w.Write([]byte("Bad Key"))
	} else if len(key) > 50 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		data := Data{}
		err := decoder.Decode(&data)

		if err != nil {
			w.Write([]byte("No value in PUT request"))
			panic(err)
		}

		s.DAL().Put(key, data.Value)
		w.Write([]byte("Put worked"))
	}
}

//GetHandler here
func (s *Store) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		w.Write([]byte("something something bad request"))
	}
	s.DAL().Get(key)
}
