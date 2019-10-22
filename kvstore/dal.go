package kvstore

import (
	"errors"
	"fmt"
)

//KVDAL is a key value data-access layer
type KVDAL struct {
	Store map[string]string // data structure
}

//Response for PUT method
const (
	ADDED   = 0
	UPDATED = 1
)

var (
	ErrKeyNotFound = errors.New("key not found")
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
	if ok {
		return value, nil
	}

	return "", fmt.Errorf("Not Found")
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
