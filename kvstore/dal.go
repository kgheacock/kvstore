package kvstore

import (
	"errors"
	"fmt"
	"sort"
)

//KVDAL is a key value data-access layer
type KVDAL struct {
	Store   map[string]string // data structure
	KeyList []string
}

//GetKeyList returns []string of all keys present in map
func (k *KVDAL) GetKeyList() ([]string, error) {
	//Clear old KeyList
	k.KeyList = nil
	//Use make for efficient memory allocation
	k.KeyList = make([]string, 0, len(k.Store))
	for key := range k.Store {
		k.KeyList = append(k.KeyList, key)
	}
	if len(k.KeyList) == 0 {
		return k.KeyList, ErrKeyListEmpty
	}
	sort.Strings(k.KeyList)
	return k.KeyList, nil
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
	fmt.Println("Putting: ", key)
	_, ok := k.Store[key]
	k.Store[key] = value
	if ok {
		return UPDATED, nil
	}
	return ADDED, nil
}

//Get function retrieves value from map if it exists
func (k *KVDAL) Get(key string) (string, error) {
	fmt.Println("Getting: ", key)
	value, ok := k.Store[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return value, nil
}

//Delete function removes key-value from map if it exists
func (k *KVDAL) Delete(key string) error {
	fmt.Println("Deleting: ", key)
	if _, ok := k.Store[key]; !ok {
		return ErrKeyNotFound
	}
	delete(k.Store, key)
	fmt.Println("Deleted ", key)
	return nil
}
