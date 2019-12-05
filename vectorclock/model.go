package vectorclock

type Store struct {
	dal DataAccessLayer
}

type DataAccessLayer interface {
	IncrementNodeClock()
	CurrentServerClock() int
	ResetVC(serverList []string)
	UpdateVC(vc2 *VectorClock)
	ReceiveEvent(incVC *VectorClock) //calls incrementlocal, then updates our VC with the incVC
}

func NewVectorClockStore(dal DataAccessLayer) *Store {
	return &Store{dal: dal}
}

func (s *Store) DAL() DataAccessLayer {
	return s.dal
}
