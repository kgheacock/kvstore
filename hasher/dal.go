package hasher
import (
	"sync"
	"hash/crc32"
	"sort"
	"strings"
)

type Hasher struct {

}

const NumVirtualNodes = 10

//Ring datastructure
type Ring struct {
	sync.Mutex
	Nodes Nodes
}

func InitRing() *Ring {
	return &Ring{Nodes: Nodes{}}
}

//List of nodes
type Nodes []*Node

//Server or Key
type Node struct {
	Id string
	IdHash uint32
}

func NewNode(id string) *Node {
	return &Node{
		Id:	id,
		IdHash: HashVal(id),
	}
}

func (r *Ring) AddKey (id string){
	r.Lock()
	defer r.Unlock()
	node := NewNode(id)
	r.Nodes = append(r.Nodes, node)
	sort.Sort(r.Nodes)
}

func (r *Ring) AddServer (id string) {
	r.Lock()
	defer r.Unlock()
	for i := 0; i < NumVirtualNodes; i++ {
		fullid := id + ":" + i
		node := NewNode(id)
		r.Nodes = append(r.Nodes, node)
	}
	sort.Sort(r.Nodes)
}

func (r *Ring) FindKey (id string) {
	location := r.search(id)
	if location >= r.Nodes.Len() {
		location = 0
	}
	serverip := strings.Split(r.Nodes[location].id, ":")
	return serverip
}

//Takes in string, returns 32bit hash
func (key string) HashVal (uint32) {
	return crc32.ChecksumIEEE([]byte(key))
}

//Return the ip of a server by passing in the key
func (h *Hasher) GetServerByKey (key string) (string, error) {
	return "localhost" , nil //temp
}
