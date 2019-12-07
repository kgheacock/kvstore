package kvstore

import (
	"sort"

	"github.com/colbyleiske/cse138_assignment2/config"
	"github.com/colbyleiske/cse138_assignment2/hasher"
	"github.com/colbyleiske/cse138_assignment2/shard"
	"github.com/colbyleiske/cse138_assignment2/vectorclock"
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
	Get(key string) (string, error)
	Put(key string, value string) (int, error)
	KeyList() []string
	GetKeyCount() int
}

func makeShards(serverList []string, replFactor int) map[int]*shard.Shard {
	sort.Strings(serverList)
	shardList := make(map[int]*shard.Shard)
	numberOfShards := len(serverList) / replFactor
	for i := 0; i < numberOfShards; i++ {
		shardMembers := serverList[i*replFactor : i*replFactor+replFactor]
		shard := shard.Shard{
			ID:          string(i),
			Nodes:       shardMembers,
			VectorClock: vectorclock.NewVectorClock(shardMembers, config.Config.Address),
		}
		shardList[i] = &shard
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

//Holds incoming PUT request body
type Data struct {
	Value string `json:"value"`
}

type ResponseMessage struct {
	Error         string                   `json:"error,omitempty"`
	Message       string                   `json:"message,omitempty"`
	Value         string                   `json:"value,omitempty"`
	Address       string                   `json:"address,omitempty"`
	CausalContext *vectorclock.VectorClock `json:"causal-context"`
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
