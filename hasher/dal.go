package hasher

import (
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const NumVirtualNodes = 20

var KeyNotFound error = errors.New("key not found")

type Ring struct {
	sync.Mutex
	Nodes   Nodes
	Numkeys int
}

func NewRing() *Ring {
	return &Ring{Nodes: Nodes{}, Numkeys: 0}
}

//List of nodes
type Nodes []*Node

func (n Nodes) Len() int           { return len(n) }
func (n Nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Nodes) Less(i, j int) bool { return n[i].IdHash < n[j].IdHash }

//Server or Key
type Node struct {
	Id       string
	IdHash   uint32
	Server   *Node
	IsServer bool
}

func (r *Ring) NewNode(id string, nodetype bool) *Node {
	return &Node{
		Id:       id,
		IdHash:   r.HashVal(id),
		Server:   nil,
		IsServer: nodetype,
	}
}

func (r *Ring) AddKey(id string) {
	r.Lock()
	defer r.Unlock()
	node := r.NewNode(id, false)
	r.Nodes = append(r.Nodes, node)
	sort.Sort(r.Nodes)
	r.Numkeys++
}

func (r *Ring) AddServer(id string) {
	r.Lock()
	defer r.Unlock()
	for i := 0; i < NumVirtualNodes; i++ {
		fullid := id + ":" + strconv.Itoa(i)
		node := r.NewNode(fullid, true)
		r.Nodes = append(r.Nodes, node)
	}
	sort.Sort(r.Nodes)
}

func (r *Ring) GetKeyNode(id string) (*Node, error) {
	boolfn := func(i int) bool {
		return r.Nodes[i].IdHash >= r.HashVal(id)
	}
	location := sort.Search(r.Nodes.Len(), boolfn)
	if location >= r.Nodes.Len() {
		location = 0
	}
	if location < r.Nodes.Len() && r.Nodes[location].Id == id {
		return r.Nodes[location], nil
	} else {
		return &Node{}, KeyNotFound
	}
}

func (r *Ring) ReShard() {
	for i := 0; i < r.Nodes.Len(); i++ {
		if !r.Nodes[i].IsServer {
			for j := i; j < r.Nodes.Len(); j++ {
				//If we reach end and didnt find server
				if j == r.Nodes.Len()-1 {
					j = 0
				}
				if r.Nodes[j].IsServer {
					r.Nodes[i].Server = r.Nodes[j]
					break
				}
			}
		}
	}
}

func (r *Ring) GetNumOfKeys() int { return r.Numkeys }

func (r *Ring) printRing() {
	for i := 0; i < r.Nodes.Len(); i++ {
		fmt.Println(r.Nodes[i].Id)
	}
}

//Takes in string, returns 32bit hash
func (r *Ring) HashVal(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

//Return the ip of a server by passing in the key
func (r *Ring) GetServerByKey(key string) (string, error) {
	node, err := r.GetKeyNode(key)
	if err != nil {
		log.Printf("Could not locate key %s\n", key)
		return "", KeyNotFound
	}
	serverip := strings.Split(node.Server.Id, ":")
	return serverip[0], nil
}

/*
//TESTING CODE
func main(){
  keys := []string{"Chris", "Brandon", "Colby", "Keith", "Alvaro", "Mackey"}

  myRing := InitRing()
  myRing.AddServer("A")
  myRing.AddServer("B")
  myRing.AddServer("C")

  for i := 0; i < len(keys); i++ {
    myRing.AddKey(keys[i])
  }

  myRing.printRing()
  myRing.ReShard()

  for i := 0; i < len(keys); i++ {
    theNodeWanted, _ := myRing.GetServerByKey(keys[i])
    fmt.Println("Key:", keys[i]+ ",","is on Server: ",theNodeWanted)
  }

  theNodeWanted, _ := myRing.GetServerByKey("DoesntExist")
  fmt.Println("Key:", "DoesntExist" + ",","is on Server: ",theNodeWanted)

}
*/
