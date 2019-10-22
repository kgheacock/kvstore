package kvstore

import (
	"crypto/sha1"
	"testing"
)

var keyList = []string{"ABCD", "FGHIJ", "LMNOP"}
var dal DataAccessLayer
var kvStore *Store

func TestInitDAL(t *testing.T) {
	dal = &KVDAL{Store: make(map[string]string)}

}
func TestInitKVStore(t *testing.T) {
	kvStore = NewStore(dal)
	if kvStore == nil {
		t.Error("kvstore was not created")
	}
}
func TestAddKVStore(t *testing.T) {
	TestInitDAL(t)
	hasher := sha1.New()
	for _, key := range keyList {
		hasher.Reset()
		hasher.Write([]byte(key))
		_, err := dal.Put(key, string(hasher.Sum(nil)))
		if err != nil {
			t.Error(err)
		}
	}
}
func TestGetKVStore(t *testing.T) {
	TestInitDAL(t)
	hasher := sha1.New()
	for _, key := range keyList {
		hasher.Reset()
		val, err := dal.Get(key)
		if err != nil {
			t.Error(err)
		}
		hasher.Write([]byte(key))
		if val != string(hasher.Sum(nil)) {
			t.Error("Incorrect key returned")
		}

	}
}

func TestDeleteKVStore(t *testing.T) {
	TestInitDAL(t)
	for _, key := range keyList {
		err := dal.Delete(key)
		if err != nil {
			t.Error(err)
		}
		_, err = dal.Get(key)
		if err == nil {
			t.Error("Key still in KVS")
		}
	}
}
