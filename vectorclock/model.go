package vectorclock

type Store struct {
	dal DataAccessLayer
}

type DataAccessLayer interface {
	IncrementVC()
	CurrentState() int
	UpdateVC(vc2 *VectorClock)
}

func NewVectorClockStore(dal DataAccessLayer) *Store {
	return &Store{dal: dal}
}

func (s *Store) DAL() DataAccessLayer {
	return s.dal
}
