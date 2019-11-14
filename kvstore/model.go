package kvstore

import "github.com/colbyleiske/cse138_assignment2/hasher"

type Store struct {
	dal                       DataAccessLayer
	hasher                    *hasher.Store
	ViewChangeFinishedChannel chan bool
	nodeCount                 int
}

type DataAccessLayer interface {
	Delete(key string) error
	Get(key string) (string, error)
	Put(key string, value string) (int, error)
	GetKeyList() ([]string, error)
}

func NewStore(dal DataAccessLayer, hasher *hasher.Store) *Store {
	return &Store{dal: dal, hasher: hasher}
}

func (s *Store) DAL() DataAccessLayer {
	return s.dal
}

func (s *Store) GetNodeCount() int {
	return s.nodeCount
}

func (s *Store) Hasher() hasher.Store {
	return *s.hasher
}

//Holds incoming PUT request body
type Data struct {
	Value string `json:"value"`
}

type ResponseMessage struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Value   string `json:"value,omitempty"`
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
