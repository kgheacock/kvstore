package kvstore

import (
	"fmt"
)

type KVDAL struct {
	//Some data strucutre
}

func (k *KVDAL) Delete(key string) error {
	fmt.Println("Deleting", key)
	return nil
}
