package hasher

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const NumVirtualNodes = 20

//*************** Node Functions ***************\\

//List of nodes
type Nodes []*Node

//Server
type Node struct {
	Ip     string
	IpHash uint32
}

func (r *Ring) NewNode(ip string) *Node {
	return &Node{
		Ip:     ip,
		IpHash: r.HashVal(ip),
	}
}

func (n Nodes) Len() int           { return len(n) }
func (n Nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Nodes) Less(i, j int) bool { return n[i].IpHash < n[j].IpHash }

//*************** Ring Functions ***************\\

//Nodes: Actual representation of ring
//Servers: []String of non-virtualized node IP's
type Ring struct {
	sync.Mutex
	Nodes   Nodes
	Servers Servers
}

//List of non-virtual server IP's
type Servers []string

func NewRing() *Ring {
	return &Ring{Nodes: Nodes{}, Servers: Servers{}}
}

//Adds server and virtual nodes to ring
func (r *Ring) AddServer(ip string) {
	r.Lock()
	defer r.Unlock()
	//Adds IP to list of servers
	r.Servers = append(r.Servers, ip)
	//Creates virtualized nodes for ring
	for i := 0; i < NumVirtualNodes; i++ {
		fullip := ip + ":" + strconv.Itoa(i)
		node := r.NewNode(fullip)
		r.Nodes = append(r.Nodes, node)
	}
	sort.Sort(r.Nodes)
}

//Returns list of all non-virtual server IP's on ring
func (r *Ring) GetServers() []string {
	sort.Strings(r.Servers)
	return r.Servers
}

//Print whole ring for debugging
func (r *Ring) printRing() {
	for i := 0; i < r.Nodes.Len(); i++ {
		fmt.Println(r.Nodes[i].Ip)
	}
}

//Takes in string, returns 32bit hash
func (r *Ring) HashVal(something string) uint32 {
	return crc32.ChecksumIEEE([]byte(something))
}

//Return the ip of a server by passing in the key
func (r *Ring) GetServerByKey(key string) (string, error) {
	//Required for binary search
	boolfn := func(i int) bool {
		return r.Nodes[i].IpHash >= r.HashVal(key)
	}
	location := sort.Search(r.Nodes.Len(), boolfn)
	//If key hashes to end of ring, server is beginning of ring
	if location >= r.Nodes.Len() {
		location = 0
	}
	node := r.Nodes[location]
	serverip := strings.Split(node.Ip, ":")
	return serverip[0], nil
}
