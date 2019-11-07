package hasher
import (
	"sync"
	"hash/crc32"
)
type Hasher struct {

}

//Ring datastructure
type Ring struct {
	sync.Mutex
	Nodes Nodes
}

//List of nodes
type Nodes []*Node

//Server or Key
type Node struct {
	Id string
	IdHash uint32
}

func (r *Ring) AddNode (id string) {
	r.Lock()
	defer r.Unlock()
}

func (r *Ring) AddKey (id string) {
	r.Lock()
	defer r.Unlock()

}



//Takes in string, returns 32bit hash
func (key string) HashId (uint32) {
	return crc32.ChecksumIEEE([]byte(key))
}

//Return the ip of a server by passing in the key
func (h *Hasher) GetServerByKey (key string) (string, error) {
	return "localhost" , nil //temp
}
