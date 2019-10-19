package kvstore

import (
	"errors"
	"fmt"
)

type KVDAL struct {
	//Some data strucutre
}

func (k *KVDAL) Delete(key string) error {
	fmt.Println("Deleting", key)
	return nil
}

func (k *KVDAL) Add(key, value string) error {
	return errors.New("Not Implemented")
}

func (k *KVDAL) Get(key string) (error, string) {
	return errors.New("Not Implemented"), ""
}
func (k *KVDAL) Update(key, value string) error {
	return errors.New("Not Implemented")
}
