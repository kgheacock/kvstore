package hasher

type Store struct {
	dal DataAccessLayer
}

type DataAccessLayer interface {
	GetServerByKey(key string) (int, error)
}

func NewHasher(dal DataAccessLayer) *Store {
	return &Store{dal: dal}
}
