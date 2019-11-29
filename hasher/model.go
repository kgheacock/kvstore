package hasher

type Store struct {
	dal DataAccessLayer
}

type DataAccessLayer interface {
	GetServerByKey(key string) (string, error)
	Servers() []string
	AddServer(ip string)
	GetQuorum(name string) (Quorum, error)
}
type Quorum struct {
	Name    string
	members []string
}

func (q *Quorum) AddMember(id string) {
	q.members = append(q.members, id)
}
func (q *Quorum) RemoveMember(id string) error {
	if len(q.members) == 0 {
		return errors.new("Element not found") //Is this correct?
	}
	if id == q.members[0] {
		return nil
	}
	for i, elem := range q.members {
		if elem == id {
			q.members = append(q.members[:i-1], q.members[i+1:]...)
			return nil
		}
	}
	return errors.new("Element not found")
}

func (q *Quorum) SetMembers(ids []string) {
	q.members = ids
}

func NewRingStore(dal DataAccessLayer) *Store {
	return &Store{dal: dal}
}

func (s *Store) DAL() DataAccessLayer {
	return s.dal
}
