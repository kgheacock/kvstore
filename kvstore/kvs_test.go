package kvstore

import (
	"crypto/sha1"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddHandler(t *testing.T) {
	TestInitDAL(t)
	TestInitKVStore(t)
	hasher := sha1.New()
	for _, key := range keyList {
		hasher.Reset()
		hasher.Write([]byte(key))
		requestBody, err := json.Marshal(map[string]string{
			"value": string(hasher.Sum(nil)),
		})
		if err != nil {
			t.Error(err)
		}
		req := httptest.NewRequest("PUT", "/kv-store/"+key, strings.NewReader(string(requestBody)))
		w := httptest.NewRecorder()
		kvStore.AddHandler(w, req)
		if w.Code != 200 {
			t.Errorf("Incorrect return code. Expected 200 got %d", w.Code)
		}
		var response map[string]string
		err = json.Unmarshal([]byte(w.Body.String()), &response)
		if err != nil {
			t.Error(err)
		}
		if response["message"] != "Added successfully" {
			t.Errorf("Incorrect message recieved. Recieved %s", response["message"])
		}
		if response["message"] != "Added successfully" {
			t.Errorf("Incorrect message recieved. Recieved %s", response["message"])
		}
	}
}

func TestDeleteHandler(t *testing.T) {
	TestInitDAL(t)
	TestInitKVStore(t)
	TestAddHandler(t)
	for _, key := range keyList {
		req := httptest.NewRequest("DELETE", "/kv-store/"+key, nil)
		w := httptest.NewRecorder()
		kvStore.DeleteHandler(w, req)
		if w.Code != 200 {
			t.Errorf("Incorrect return code. Expected 200 got %d", w.Code)
		}
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		if err != nil {
			t.Error(err)
		}
		if response["message"] != "Deleted successfully" {
			t.Errorf("Incorrect message recieved. Recieved %s", response["message"])
		}
	}
}

func TestGetHandler(t *testing.T) {
	TestInitDAL(t)
	TestInitKVStore(t)
	TestAddHandler(t)
	hasher := sha1.New()
	for _, key := range keyList {
		hasher.Reset()
		hasher.Write([]byte(key))
		req := httptest.NewRequest("PUT", "/kv-store/"+key, nil)
		w := httptest.NewRecorder()
		kvStore.DeleteHandler(w, req)
		if w.Code != 200 {
			t.Errorf("Incorrect return code. Expected 200 got %d", w.Code)
		}
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		if err != nil {
			t.Error(err)
		}
		if response["message"] != "Retrieved successfully" {
			t.Errorf("Incorrect message recieved. Recieved %s", response["message"])
		}
		if response["value"] != string(hasher.Sum(nil)) {
			t.Errorf("Recieved %s but expected %s", response["value"], string(hasher.Sum(nil)))
		}
	}
}
