package kvstore

type Store struct {
	dal DataAccessLayer
}

type DataAccessLayer interface {
	//Add
	Delete ()
	//Potential Update
}

func NewStore( dal DataAccessLayer) *Store {
	return &Store{dal: dal}
}

func (s *Store) DAL() DataAccessLayer {
	return s.dal
}