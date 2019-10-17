package kvstore

import (
	"encoding/json"
	"net/http"
	"./errormessages"
	"github.com/gorilla/mux"
)

//Data is the json repsonse
type Data struct {
	Value string `json:"value"`
}

//DeleteHandler here
func (s *Store) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	key, ok := vars["key"]
	errMsg := DeleteFailure{} //For failure formatting
	delMsg := DeleteSuccess{} //For success put formatting
	if !ok { //Key is not in URL
		errMsg.Exists = false
		errMsg.Error = "No key"
		errMsg.Message = "Error in DELETE"
		w.WriteHeader(http.StatusNotFound) //404
		json.NewEncoder(w).Enode(errMsg)
	} else if len(key) > 50 { //Key too long
		errMsg.Exists = false
		errMsg.Error = "Key is too long"
		errMsg.Message = "Error in DELETE"
		w.WriteHeader(http.StatusBadRequest) //400
		json.NewEncoder(w).Encode(errMsg)
	} else { //Some key present in URL
		val, err := s.DAL().Get(key)
		if err != nil { //Error
			errMsg.Exists = false
			errMsg.Error = "Key does not exist"
			errMsg.Message = "Error in DELETE"
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errMsg)
		} else { //Key exists in KVS
			delMsg.Exists = true
			delMsg.Message = "Deleted successfully"
			w.WriteHeader(http.StatusOK)
			s.DAL().Delete(key)
			json.NewEncoder(w).Encode(delMsg)
		}
	}
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

		if err == nil {
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

	if !ok {
		w.Write([]byte("Method GET not supported"))
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		val, err := s.DAL().Get(key)
		if err != nil {
			getMsg.Exists = false
			getMsg.Value = "Key does not exist"
			getMsg.Message = "Error in GET"
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
