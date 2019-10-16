package kvstore

type Store struct {
	dal DataAccessLayer
}

type DataAccessLayer interface {
	Delete(key string) error

	Put(key string, value string) error

	Get(key string) error
}

func NewStore(dal DataAccessLayer) *Store {
	return &Store{dal: dal}
}

func (s *Store) DAL() DataAccessLayer {
	return s.dal
}
