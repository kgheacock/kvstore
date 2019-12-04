package kvstore

import (
	"errors"
	"log"
	"sort"
)

type REPL_RESULT struct {
	FAIL    "FAIL"
	SUCCESS "SUCCESS"
}

//KVDAL is a key value data-access layer
type KVDAL struct {
	Store   map[string]string // data structure
	keyList []string
}

//GetKeyList returns []string of all keys present in map
func (k *KVDAL) KeyList() []string {
	//Clear old KeyList
	k.keyList = nil
	//Use make for efficient memory allocation
	k.keyList = make([]string, 0, len(k.Store))
	for key := range k.Store {
		k.keyList = append(k.keyList, key)
	}

	sort.Strings(k.keyList)
	return k.keyList
}

//Response for PUT method
const (
	ADDED   = 0
	UPDATED = 1
)

var (
	ErrKeyNotFound  = errors.New("key not found")
	ErrKeyListEmpty = errors.New("key list empty")
)

//Put function stores value into map based on key
func (k *KVDAL) Put(key string, value string) (int, error) {
	log.Println("PUT", key)
	_, ok := k.Store[key]
	k.Store[key] = value
	if ok {
		return UPDATED, nil
	}
	return ADDED, nil
}

//Get function retrieves value from map if it exists
func (k *KVDAL) Get(key string) (string, error) {
	log.Println("GET", key)
	value, ok := k.Store[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return value, nil
}

//Delete function removes key-value from map if it exists
func (k *KVDAL) Delete(key string) error {
	log.Println("DEL", key)
	if _, ok := k.Store[key]; !ok {
		log.Printf("%s not found\n", key)
		return ErrKeyNotFound
	}
	delete(k.Store, key)
	return nil
}

//GetKeyCount gets the number of keys in the map
func (k *KVDAL) GetKeyCount() int {
	return len(k.Store)
}
