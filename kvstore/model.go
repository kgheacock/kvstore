package kvstore

import (
	"sort"

	"github.com/colbyleiske/cse138_assignment2/hasher"
	"github.com/colbyleiske/cse138_assignment2/shard"
)

type Store struct {
	dal    DataAccessLayer
	hasher *hasher.Store
	state  nodeState
}
type nodeState int

const (
	NORMAL nodeState = iota + 1
	RECIEVED_EXTERNAL_RESHARD
	//Lock Dict to external requests
	RECIEVED_INTERNAL_RESHARD
	PREPARE_FOR_RESHARD
	TRANSFER_KEYS
	FINISHED_TRANSFER
	WAITING_FOR_ACK
	//Release lock
	PROCESS_BACKLOG
)

type DataAccessLayer interface {
	Delete(key string) error
	Get(key string) (StoredValue, error)
	Put(key string, value StoredValue) int
	KeyList() []string
	GetKeyCount() int
	MapKeyToClock() (keyClock map[string]int)
	IncrementClock(key string) int //returns new clock value
	SetClock(key string, newClock int) 
}

func makeShards(serverList []string, replFactor int) map[int]*shard.Shard {
	//warning - does not set our current shard ID
	sort.Strings(serverList)
	shardList := make(map[int]*shard.Shard)
	numberOfShards := len(serverList) / replFactor
	for i := 0; i < numberOfShards; i++ {
		shardMembers := serverList[i*replFactor : i*replFactor+replFactor]
		shard := shard.Shard{
			ID:    string(i),
			Nodes: shardMembers,
		}
		shardList[i] = &shard
		// if config.Contains(servers[0+(replFactorNum*i):replFactorNum+(replFactorNum*i)], addr) {
		// 	config.Config.CurrentShardID = i
		// }
	}
	return shardList

}
func NewStore(dal DataAccessLayer, hasher *hasher.Store) *Store {
	return &Store{dal: dal, hasher: hasher, state: NORMAL}
}

func (s *Store) DAL() DataAccessLayer {
	return s.dal
}

func (s *Store) Hasher() hasher.Store {
	return *s.hasher
}

func (s *Store) State() nodeState {
	return s.state
}

type StoredValue struct {
	value        string
	lamportclock int
}

//Holds incoming PUT request body
type Data struct {
	Value string `json:"value"`
}

type ResponseMessage struct {
	Error         string         `json:"error,omitempty"`
	Message       string         `json:"message,omitempty"`
	Value         string         `json:"value,omitempty"`
	Address       string         `json:"address,omitempty"`
	CausalContext map[string]int `json:"causal-context"`
}

type DeleteResponse struct {
	ResponseMessage
	Exists bool `json:"doesExist"`
}

type PutResponse struct {
	ResponseMessage
	Replaced bool `json:"replaced"`
}

type GetResponse struct {
	ResponseMessage
	Exists bool `json:"doesExist"`
}

type GetKeyCountRepsponse struct {
	Message  string `json:"message"`
	KeyCount int    `json:"key-count"`
}

type ViewChangeRequest struct {
	View       []string `json:"view"`
	ReplFactor int      `json:"repl-factor"`
}

type NodeStatus struct {
	IP       string
	KeyCount int
	ShardID  int
}
