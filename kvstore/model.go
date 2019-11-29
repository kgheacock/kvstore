package kvstore

import (
	"github.com/colbyleiske/cse138_assignment2/hasher"
)

type Store struct {
	dal                       DataAccessLayer
	hasher                    *hasher.Store
	state                     nodeState
	ViewChangeFinishedChannel chan bool
}
type shard struct {
	Address  string `json:"address,omitempty"`
	KeyCount int    `json:"key-count"`
}


type DataAccessLayer interface {
	Delete(key string) error
	Get(key string) (string, error)
	Put(key string, value string) (int, error)
	KeyList() []string
	GetKeyCount() int
}

func NewStore(dal DataAccessLayer, hasher *hasher.Store) *Store {
	return &Store{dal: dal, hasher: hasher, state: NORMAL, ViewChangeFinishedChannel: make(chan bool, 1)}
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
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Value   string `json:"value,omitempty"`
	Address string `json:"address,omitempty"`
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

type ExternalViewChangeRequest struct {
	View string[] `json:"view"`
	ReplFactor int `json:"repl-factor"`
}
type InternalViewChangeRequest struct{
	NamedView map[string]string[] `json:"named-view"`
	ReplFactor int `json:"repl-factor"`
}