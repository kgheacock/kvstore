package hasher

type Store struct {
	dal DataAccessLayer
}

type DataAccessLayer interface {
	GetServerByKey(key string) (string, error)
}

func NewHasher(dal DataAccessLayer) *Store {
	return &Store{dal: dal}
}
