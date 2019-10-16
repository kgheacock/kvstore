package kvstore

import (
	"fmt"
)

//KVDAL is a key value data-access layer
type KVDAL struct {
	Store map[string]string
}

//Put here
func (k *KVDAL) Put(key string, value string) error {
	fmt.Println("Putting", key)
	k.Store[key] = value
	return nil
}

//Get here
func (k *KVDAL) Get(key string) error {
	fmt.Println("Getting", key)
	value, ok := k.Store[key]
	if ok {
		fmt.Println("Got ", value)
	} else {
		fmt.Printf("Not Value")
	}
	return nil
}

//Delete here
func (k *KVDAL) Delete(key string) error {
	fmt.Println("Deleting", key)
	_, ok := k.Store[key]
	if ok {
		delete(k.Store, key)
		fmt.Println("Deleted ", key)
	} else {
		fmt.Printf("No value")
	}
	return nil
}
