package kvstore

import (
	"fmt"
)

//KVDAL is a key value data-access layer
type KVDAL struct {
	Store map[string]string // data structure
}

//Put here
func (k *KVDAL) Put(key string, value string) (string, error) {
	fmt.Println("Putting", key)
	_, ok := k.Store[key]
	k.Store[key] = value
	if ok {
		return "replaced", nil
	}
	return "added", nil
}

//Get here
func (k *KVDAL) Get(key string) (string, error) {
	fmt.Println("Getting", key)
	value, ok := k.Store[key]
	if ok {
		fmt.Println("Got ", value)
	} else {
		fmt.Printf("Not Value")
	}
	return value, nil
}

//Delete here
func (k *KVDAL) Delete(key string) error {
	fmt.Println("Deleting", key)
	_, ok := k.Store[key]
	if ok {
		delete(k.Store, key)
		fmt.Println("Deleted ", key)
		return nil
	}
	return fmt.Errorf("Key Not Valid")
}
