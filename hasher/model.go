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

func (q *Quorum) AddMember(ip string) {
	q.members = append(q.members, ip)
}
func (q *Quorum) RemoveMember(ip string) error{
	if len(q.members) == 0{
		return Errors.new("Element not found") //Is this correct?
	}
	if elem == q.members[0]{
		return nil
	}
	for i,elem := q.members{
		if(elem == ip){
			q.members = append(q.members[:i-1],q.members[i+1:])
			return nil
		}
	}
	return Errors.new("Element not found")
}

func (q *Quorum) SetMembers(ips [] string){
	q.members = ips
}

func NewRingStore(dal DataAccessLayer) *Store {
	return &Store{dal: dal}
}

func (s *Store) DAL() DataAccessLayer {
	return s.dal
}
