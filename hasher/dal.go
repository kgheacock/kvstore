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

type HasherError struct {
  Key string
  Message string
}

func (e HasherError) Error() string {
  return fmt.Sprintf("%v: %v", e.Key, e.Message)
}

func KeyNotFound(key string) error {
  return HasherError{
    key,
    "could not be located.",
  }
}

//Ring datastructure
type Ring struct {
  sync.Mutex
  Nodes Nodes
  Numkeys int
}

func InitRing() *Ring {
  return &Ring{Nodes: Nodes{}, Numkeys: 0}
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
  Server *Node
  IsServer bool
}

func NewNode(id string, nodetype bool) *Node {
  return &Node{
    Id:  id,
    IdHash: HashVal(id),
    Server: nil,
    IsServer: nodetype,
  }
}

func (r *Ring) AddKey (id string){
  r.Lock()
  defer r.Unlock()
  node := NewNode(id, false)
  r.Nodes = append(r.Nodes, node)
  sort.Sort(r.Nodes)
  r.Numkeys++
}

func (r *Ring) AddServer (id string) {
  r.Lock()
  defer r.Unlock()
  for i := 0; i < NumVirtualNodes; i++ {
    fullid := id + ":" + strconv.Itoa(i)
    node := NewNode(fullid, true)
    r.Nodes = append(r.Nodes, node)
  }
  sort.Sort(r.Nodes)
}

func (r *Ring) GetKeyNode (id string) (*Node, error) {
  boolfn := func(i int) bool {
    return r.Nodes[i].IdHash >= HashVal(id)
  }
  location := sort.Search(r.Nodes.Len(), boolfn)
  if location >= r.Nodes.Len() {
    location = 0
  }
  if location == -1 { //***THIS DOES NOT CURRENTLY WORK***
    return &Node{}, KeyNotFound("")
  }
  return r.Nodes[location], nil
}

func (r *Ring) ReShard () {
  for i := 0; i < r.Nodes.Len(); i++ {
    if r.Nodes[i].IsServer == false {
      for j := i; j < r.Nodes.Len(); j++ {
        //If we reach end and didnt find server
        if j == r.Nodes.Len() - 1 {
          j = 0
        }
        if r.Nodes[j].IsServer == true {
          r.Nodes[i].Server = r.Nodes[j]
          break
        }
      }
    }
  }
}

func (r *Ring) GetNumOfKeys () int {return r.Numkeys}

func (r *Ring) printRing () {
  for i := 0; i < r.Nodes.Len(); i++ {
    fmt.Println(r.Nodes[i].Id)
  }
}

//Takes in string, returns 32bit hash
func HashVal (key string) uint32 {
  return crc32.ChecksumIEEE([]byte(key))
}

//Return the ip of a server by passing in the key
func (r *Ring) GetServerByKey (key string) (string, error) {
  node, err := r.GetKeyNode(key)
  if err != nil {
    return "could not locate key.", KeyNotFound(key)
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
  
  //This breaks everything
  theNodeWanted, _ := myRing.GetServerByKey("DoesntExist")
  fmt.Println("Key:", "DoesntExist" + ",","is on Server: ",theNodeWanted)
  
}
*/
