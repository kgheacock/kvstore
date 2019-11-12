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

//*************** Node Functions ***************\\

//List of nodes
type Nodes []*Node

//Server or Key
type Node struct {
	Id       string
	IdHash   uint32
	Server   *Node
	IsServer bool
}

func (n Nodes) Len() int           { return len(n) }
func (n Nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Nodes) Less(i, j int) bool { return n[i].IdHash < n[j].IdHash }

//*************** Ring Functions ***************\\

type Ring struct {
	sync.Mutex
	Nodes   Nodes
	Numkeys int
}

func NewRing() *Ring {
	return &Ring{Nodes: Nodes{}, Numkeys: 0}
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

func (r *Ring) AddServer(ip string) {
	r.Lock()
	defer r.Unlock()
	for i := 0; i < NumVirtualNodes; i++ {
		fullip := ip + ":" + strconv.Itoa(i)
		node := r.NewNode(fullip, true)
		r.Nodes = append(r.Nodes, node)
	}
	sort.Sort(r.Nodes)
}

//Used to binary search for node in ring
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

//Get total keys on ring
func (r *Ring) GetNumOfKeys() int { return r.Numkeys }

//Print whole ring for debugging
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

//Reassign keys to servers
func (r *Ring) ReShard() {
	//Empty slice for resharding
	ServersAndKeys = nil
	//Empty ring
	if r.GetNumOfKeys() == 0 {
		return
	}
	r.Lock()
	defer r.Unlock()
	for i := 0; i < r.Nodes.Len(); i++ {
		if !r.Nodes[i].IsServer {
			for j := i; j < r.Nodes.Len(); j++ {
				//If we reach end and didnt find server
				if j == r.Nodes.Len()-1 {
					j = 0
				}
				if r.Nodes[j].IsServer {
					r.Nodes[i].Server = r.Nodes[j]
					KeysPerNode(r.Nodes[i], r.Nodes[j])
					break
				}
			}
		}
	}
}

//*************** KeyCount Functions ***************\\

//Holds struct containing server name and its keys
type KeyCount struct {
	ServerName string
	Nodes      Nodes
}

//Defines slice of KeyCount structs
var ServersAndKeys []KeyCount

func NewKeyCount(name string) *KeyCount {
	return &KeyCount{
		ServerName: name,
		Nodes:      Nodes{},
	}
}

//Called during resharding, creates list of KeyCount structs
//Used to map data to respective servers outside this library
func KeysPerNode(key, server *Node) {
	serverIp := strings.Split(server.Id, ":")
	for entry := range ServersAndKeys {
		if ServersAndKeys[entry].ServerName == serverIp[0] {
			ServersAndKeys[entry].Nodes = append(ServersAndKeys[entry].Nodes, key)
			return
		}
	}
	//If no entry exists for that name
	newEntry := NewKeyCount(serverIp[0])
	newEntry.Nodes = append(newEntry.Nodes, key)
	ServersAndKeys = append(ServersAndKeys, *newEntry)
	return
}

func GetServersAndKeys() []KeyCount { return ServersAndKeys }

/*
Fix for loops
Fix name scheme for servers id methods
*/
