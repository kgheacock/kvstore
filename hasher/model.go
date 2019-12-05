package hasher

import "github.com/colbyleiske/cse138_assignment2/shard"

type Store struct {
	dal DataAccessLayer
}

type DataAccessLayer interface {
	GetServerByKey(key string) (string, error)
	Shards() Shards
	AddShard(newShard *shard.Shard)
	RemoveShard(shard *shard.Shard)
}

func NewRingStore(dal DataAccessLayer) *Store {
	return &Store{dal: dal}
}

func (s *Store) DAL() DataAccessLayer {
	return s.dal
}
