package router

import (
	"github.com/colbyleiske/cse138_assignment2/hasher"
	"github.com/colbyleiske/cse138_assignment2/kvstore"
)

type Store struct {
	hasher  *hasher.Store
	kvstore *kvstore.Store
}

func NewStore(hasher *hasher.Store, kvstore *kvstore.Store) *Store {
	return &Store{hasher: hasher, kvstore: kvstore}
}
