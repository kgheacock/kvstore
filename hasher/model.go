package hasher

type Store struct {
	dal DataAccessLayer
}

type DataAccessLayer interface {
	Hash(key string) (int, error)
	GetNodeCount() (int) //temp - could change this to global config... Or make it a wrapper for global config value
}

func NewHasher(dal DataAccessLayer) *Store {
	return &Store{dal: dal}
}