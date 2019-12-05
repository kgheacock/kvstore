package hasher

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"sort"
	"strconv"

	"github.com/colbyleiske/cse138_assignment2/shard"
)

const numVirtualNodes = 20

//*************** Node Functions ***************\\

type nodes []*node

//node is a server
type node struct {
	Shard  shard.Shard
	IPHash uint32
}

func (r *Ring) newNode(shard shard.Shard) *node {
	return &node{
		Shard:  shard,
		IPHash: r.hashVal(shard.ID),
	}
}

func (n nodes) Len() int           { return len(n) }
func (n nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n nodes) Less(i, j int) bool { return n[i].IPHash < n[j].IPHash }

type Shards []*shard.Shard

func (n Shards) Len() int           { return len(n) }
func (n Shards) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Shards) Less(i, j int) bool { return n[i].ID < n[j].ID }

//*************** Ring Functions ***************\\

//Ring nodes: Actual representation of ring by virtual nodes
//Ring servers: []String of non-virtualized node IP's
type Ring struct {
	nodes  nodes
	shards Shards
}

//NewRing creates Ring object
func NewRing() *Ring {
	return &Ring{nodes: nodes{}, shards: Shards{}}
}

//AddServer adds server and virtual nodes to ring
func (r *Ring) AddShard(newShard *shard.Shard) {
	//Adds IP to list of servers
	r.shards = append(r.shards, newShard)
	newVirNodes := make(nodes, 0, numVirtualNodes)
	//Creates virtualized nodes for ring
	for i := 0; i < numVirtualNodes; i++ {
		virtualShard := shard.Shard{Nodes: newShard.Nodes, ID: newShard.ID + "$" + strconv.Itoa(i), VectorClock: newShard.VectorClock}
		node := r.newNode(virtualShard)
		newVirNodes = append(newVirNodes, node)
	}
	r.nodes = append(r.nodes, newVirNodes...)
	sort.Sort(r.nodes)
}

//RemoveServer removes a server from the ring, and changes server list
func (r *Ring) RemoveShard(shard *shard.Shard) {
	location := 0
	for i, s := range r.Shards() {
		if s.ID == shard.ID {
			location = i
		}
	}

	newShards := append(r.Shards()[:location], r.Shards()[location+1:]...)
	//Clear nodes entirely and remake ring
	r.nodes = nil
	for _, item := range newShards {
		r.AddShard(item)
	}
	r.shards = newShards
}

//Servers returns list of all non-virtual server IP's on ring
func (r *Ring) Shards() Shards {
	sort.Sort(r.shards)
	return r.shards
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
	nodeIP := node.Shard.Nodes[rand.Intn(len(node.Shard.Nodes))] // for now - pick a random replica :)
	return nodeIP, nil
}

func (r *Ring) PrintRing() {
	for i := 0; i < r.nodes.Len(); i++ {
		fmt.Println(r.nodes[i].Shard.ID)
	}
}

func (r *Ring) hashVal(val string) uint32 {
	return crc32.ChecksumIEEE([]byte(val))
}
