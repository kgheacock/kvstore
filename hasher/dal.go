package hasher

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
)

const numVirtualNodes = 20

//*************** Node Functions ***************\\

type nodes []*node

//node is a server
type node struct {
	IP     string
	IPHash uint32
}

func (r *Ring) newNode(ip string) *node {
	return &node{
		IP:     ip,
		IPHash: r.hashVal(ip),
	}
}

func (n nodes) Len() int           { return len(n) }
func (n nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n nodes) Less(i, j int) bool { return n[i].IPHash < n[j].IPHash }

//*************** Ring Functions ***************\\

//Ring nodes: Actual representation of ring by virtual nodes
//Ring servers: []String of non-virtualized node IP's
type Ring struct {
	nodes   nodes
	servers []string
}

//NewRing creates Ring object
func NewRing() *Ring {
	return &Ring{nodes: nodes{}, servers: []string{}}
}

//AddServer adds server and virtual nodes to ring
func (r *Ring) AddServer(ip string) {
	//Adds IP to list of servers
	r.servers = append(r.servers, ip)
	newVirNodes := make(nodes, 0, numVirtualNodes)
	//Creates virtualized nodes for ring
	for i := 0; i < numVirtualNodes; i++ {
		virtualIP := ip + "$" + strconv.Itoa(i)
		node := r.newNode(virtualIP)
		newVirNodes = append(newVirNodes, node)
	}
	r.nodes = append(r.nodes, newVirNodes...)
	sort.Sort(r.nodes)
}

//RemoveServer removes a server from the ring, and changes server list
func (r *Ring) RemoveServer(ip string) {
	location := sort.SearchStrings(r.Servers(), ip)
	newServers := append(r.Servers()[:location], r.Servers()[location+1:]...)
	//Clear nodes entirely and remake ring
	r.nodes = nil
	for _, item := range newServers {
		r.AddServer(item)
	}
	r.servers = newServers
}

//Servers returns list of all non-virtual server IP's on ring
func (r *Ring) Servers() []string {
	sort.Strings(r.servers)
	return r.servers
}

//GetServerByKey returns the IP of a server by passing in the key
func (r *Ring) GetServerByKey(key string) (string, error) {
	//Required for binary search
	boolfn := func(i int) bool {
		return r.nodes[i].IPHash >= r.hashVal(key)
	}
	location := sort.Search(r.nodes.Len(), boolfn)
	//If key hashes to end of ring, server is beginning of ring
	if location >= r.nodes.Len() {
		location = 0
	}
	node := r.nodes[location]
	serverip := strings.Split(node.IP, "$")
	return serverip[0], nil
}

func (r *Ring) printRing() {
	for i := 0; i < r.nodes.Len(); i++ {
		fmt.Println(r.nodes[i].IP)
	}
}

func (r *Ring) hashVal(val string) uint32 {
	return crc32.ChecksumIEEE([]byte(val))
}
