package hasher

import (
  "sync"
  "hash/crc32"
  "sort"
  "strings"
  "fmt"
  "strconv"
)

const NumVirtualNodes = 20

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

func (n Nodes) Len() int           { return len(n) }
func (n Nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Nodes) Less(i, j int) bool { return n[i].IdHash < n[j].IdHash }

//Server or Key
type Node struct {
  Id string
  IdHash uint32
}

func NewNode(id string) *Node {
  return &Node{
    Id:  id,
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
    fullid := id + ":" + strconv.Itoa(i)
    node := NewNode(fullid)
    r.Nodes = append(r.Nodes, node)
  }
  sort.Sort(r.Nodes)
}

func (r *Ring) FindKey (id string) string {
  //location := sort.Search(r.Nodes.Len(), id)
  boolfn := func(i int) bool {
    return r.Nodes[i].IdHash >= HashVal(id)
  }
  //if location >= r.Nodes.Len() {
  //  location = 0
  //}
  location := sort.Search(r.Nodes.Len(), boolfn)
  serverip := strings.Split(r.Nodes[location].Id, ":")
  return serverip[0]
}

func (r *Ring) printRing() {
  for i := 0; i < r.Nodes.Len(); i++ {
    fmt.Println(r.Nodes[i].Id)
  }
}

//func (r *Ring) search(id string) int {
//  searchfn := func(i int) bool {
//    return r.Nodes[i].HashId >= hashId(id)
//  }
//
//  return sort.Search(r.Nodes.Len(), searchfn)
//}

//Takes in string, returns 32bit hash
func HashVal(key string) uint32 {
  return crc32.ChecksumIEEE([]byte(key))
}

//Return the ip of a server by passing in the key
//func (h *Hasher) GetServerByKey (key string) (string, error) {
//  return "localhost" , nil //temp
//}