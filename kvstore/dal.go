package kvstore

import (
	"fmt"
)

type KVDAL struct {
	//Some data strucutre
}

func (k *KVDAL) Delete() {
	fmt.Println("DELETE")
}