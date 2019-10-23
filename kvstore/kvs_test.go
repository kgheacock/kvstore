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
		bodyStruct := struct {
			Value string `json:"value"`
		} {
			Value: string(hasher.Sum(nil)),
		}
		requestBody, err := json.Marshal(bodyStruct)
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
func TestAddHandlerExists(t *testing.T) {
	TestInitDAL(t)
	TestInitKVStore(t)
	hasher := sha1.New()
	hasher.Write([]byte(keyList[0]))
	requestBody, err := json.Marshal(map[string]string{
		"value": string(hasher.Sum(nil)),
	})
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest("PUT", "/kv-store/"+keyList[0], strings.NewReader(string(requestBody)))
	w := httptest.NewRecorder()
	kvStore.AddHandler(w, req)
	if w.Code != 200 {
		t.Errorf("Incorrect return code. Expected 200 got %d", w.Code)
	}
	requestBodyDup, _ := json.Marshal(map[string]string{
		"value": "ABC",
	})
	reqDup := httptest.NewRequest("PUT", "/kv-store/"+keyList[0], strings.NewReader(string(requestBodyDup)))
	wDup := httptest.NewRecorder()
	kvStore.AddHandler(wDup, reqDup)
	if wDup.Code != 200 {
		t.Errorf("Incorrect return code. Expected 200 got %d", wDup.Code)
	}
	var response map[string]string
	err = json.Unmarshal([]byte(wDup.Body.String()), &response)
	if response["replaced"] == "True" {
		t.Errorf("Duplicate key should cause value to update. Response %s", response["replaced"])
	}

}
// func TestAddHandlerTooLong(t *testing.T) {
// 	TestInitDAL(t)
// 	TestInitKVStore(t)
// 	hasher := sha1.New()
// 	hasher.Write([]byte(keyList[0]))
// 	requestBody, err := json.Marshal(map[string]string{
// 		"value": string(hasher.Sum(nil)),
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	req := httptest.NewRequest("PUT", "/kv-store/"+string([0...99]), strings.NewReader(string(requestBody)))
// 	w := httptest.NewRecorder()
// 	kvStore.AddHandler(w, req)
// 	if w.Code != 200 {
// 		t.Errorf("Incorrect return code. Expected 200 got %d", w.Code)
// 	}

// }

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