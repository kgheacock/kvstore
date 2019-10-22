package kvstore

import (
	"crypto/sha1"
	"testing"
)

var keyList = []string{"ABCD", "FGHIJ", "LMNOP"}
var dal KVDAL
var kvStore *Store

func TestInitDAL(t *testing.T) {
	dal = KVDAL{}
}
func TestInitKVStore(t *testing.T) {
	kvStore = NewStore(&dal)
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
		err := dal.Add(key, string(hasher.Sum(nil)))
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
		err, val := dal.Get(key)
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
		err1, _ := dal.Get(key)
		if err1 == nil {
			t.Error("Key still in KVS")
		}
	}
}
