package hasher

type Store struct {
	dal DataAccessLayer
}

type DataAccessLayer interface {
	ServerOfKey(key string) (string, error)
	Servers() []string
	AddServer(ip string)
}

func NewRingStore(dal DataAccessLayer) *Store {
	return &Store{dal: dal}
}

func (s *Store) DAL() DataAccessLayer {
	return s.dal
}
