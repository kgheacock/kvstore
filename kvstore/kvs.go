package kvstore

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

//DeleteHandler
func (s *Store) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	returnMsg := ResponseMessage{}
	//Key is not in URL
	if !ok {
		//exists.Exists is neccessary because anonymous function "exists"
		//contains the value Exists. This is required due to use of
		//omitempty in our JSON objects
		returnMsg.exists.Exists = false
		returnMsg.Error = "No key"
		returnMsg.Message = "Error in DELETE"
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(returnMsg)
	//Key present in URL but is too long
	} else if len(key) > 50 {
		returnMsg.exists.Exists = false
		returnMsg.Error = "Key is too long"
		returnMsg.Message = "Error in DELETE"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(returnMsg)
	//Key present in URL
	} else {
		_, err := s.DAL().Get(key)
		if err != nil {
			returnMsg.exists.Exists = false
			returnMsg.Error = "Key does not exist"
			returnMsg.Message = "Error in DELETE"
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(returnMsg)
		//Key exists in KVS
		} else {
			returnMsg.exists.Exists = true
			returnMsg.Message = "Deleted successfully"
			w.WriteHeader(http.StatusOK)
			s.DAL().Delete(key)
			json.NewEncoder(w).Encode(returnMsg)
		}
	}
}

//PutHandler puts a key into the store
func (s *Store) PutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	key, ok := vars["key"]
	returnMsg := ResponseMessage{}

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No key"))
	} else if len(key) > 50 {
		returnMsg.Error = "Key is too long"
		returnMsg.Message = "Error in PUT"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(returnMsg)
	} else {
		data := Data{}
		err := decoder.Decode(&data)

		if err == nil {
			if data.Value == "" {
				returnMsg.Error = "Value is missing"
				returnMsg.Message = "Error in PUT"
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(returnMsg)
			} else {
				valStr, err := s.DAL().Put(key, data.Value)

				if err != nil {
					panic(err)
				} else {
					if valStr == "added" {
						returnMsg.Message = "Added successfully"
						//replaced.Replaced is neccessary because anonymous function "replaced"
						//contains the value Replaced. This is required due to use of
						//omitempty in our JSON objects
						returnMsg.replaced.Replaced = false
						w.WriteHeader(http.StatusCreated)
					} else if valStr == "replaced" {
						returnMsg.Message = "Updated successfully"
						returnMsg.replaced.Replaced = true
						w.WriteHeader(http.StatusOK)
					}
					json.NewEncoder(w).Encode(returnMsg)
				}
			}
		}
	}
}

//GetHandler
func (s *Store) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	returnMsg := ResponseMessage{}

	if !ok {
		w.Write([]byte("Method GET not supported"))
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		val, err := s.DAL().Get(key)
		if err != nil {
			//exists.Exists is neccessary because anonymous function "exists"
			//contains the value Exists. This is required due to use of
			//omitempty in our JSON objects
			returnMsg.exists.Exists = false
			returnMsg.Error = "Key does not exist"
			returnMsg.Message = "Error in GET"
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(returnMsg)
		} else {
			returnMsg.exists.Exists = true
			returnMsg.Value = val
			returnMsg.Message = "Retrieved successfully"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(returnMsg)
		}
	}
}
