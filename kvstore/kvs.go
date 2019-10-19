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
	errMsg := DeleteFailure{} //For failure formatting
	delMsg := DeleteSuccess{} //For success put formatting
	if !ok {                  //Key is not in URL
		errMsg.Exists = false
		errMsg.Error = "No key"
		errMsg.Message = "Error in DELETE"
		w.WriteHeader(http.StatusNotFound) //404
		json.NewEncoder(w).Encode(errMsg)
	} else if len(key) > 50 { //Key too long
		errMsg.Exists = false
		errMsg.Error = "Key is too long"
		errMsg.Message = "Error in DELETE"
		w.WriteHeader(http.StatusBadRequest) //400
		json.NewEncoder(w).Encode(errMsg)
	} else { //Some key present in URL
		_, err := s.DAL().Get(key)
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
				putResp, err := s.DAL().Put(key, data.Value)

				if err != nil {
					panic(err)
				} else {
					if putResp == ADDED {
						putMsg.Message = "Added successfully"
						putMsg.Replaced = false
						w.WriteHeader(http.StatusCreated)
					} else if putResp == UPDATED {
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
			json.NewEncoder(w).Encode(errMsg)
		} else {
			getMsg.Exists = true
			getMsg.Value = val
			getMsg.Message = "Retrieved successfully"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(getMsg)
		}
	}
}
