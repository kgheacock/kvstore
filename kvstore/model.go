package kvstore

//Store ...
type Store struct {
	dal DataAccessLayer
}

//DataAccessLayer interface
type DataAccessLayer interface {
	Delete(key string) error

	Put(key string, value string) (int, error)

	Get(key string) (string, error)
}

//NewStore creates a store
func NewStore(dal DataAccessLayer) *Store {
	return &Store{dal: dal}
}

//DAL data access layer
func (s *Store) DAL() DataAccessLayer {
	return s.dal
}
