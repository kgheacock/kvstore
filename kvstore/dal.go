package kvstore

import (
	"errors"
	"log"
	"sort"
)

//KVDAL is a key value data-access layer

type KVDAL struct {
	Store   map[string]StoredValue // data structure
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
func (k *KVDAL) Put(key string, value StoredValue) int {
	log.Println("PUT", key)
	_, ok := k.Store[key]
	k.Store[key] = value
	if ok {
		return UPDATED
	}
	return ADDED
}

//Get function retrieves value from map if it exists
func (k *KVDAL) Get(key string) (StoredValue, error) {
	log.Println("GET", key)
	value, ok := k.Store[key]
	if !ok {
		return StoredValue{}, ErrKeyNotFound
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

func (k *KVDAL) IncrementClock(key string) int {
	if val, ok := k.Store[key]; ok {
		val.lamportclock++
		k.Store[key] = val
		return val.lamportclock
	}
	// if we do not have a clock value - this means we are adding it for t he first time... start with 1 on context
	return 1
}

func (k *KVDAL) SetClock(key string, newClock int) {

}

//GetKeyCount gets the number of keys in the map
func (k *KVDAL) GetKeyCount() int {
	return len(k.Store)
}

func (k *KVDAL) MapKeyToClock() map[string]int {
	keyClock := make(map[string]int)
	for k, v := range k.Store {
		keyClock[k] = v.lamportclock
	}
	return keyClock
}
