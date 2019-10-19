package kvstore

//Store ...
type Store struct {
	dal DataAccessLayer
}

type Data struct {
	Value string
}

//Error and Success message
type ResponseMessage struct {
	Exists   *bool  `json:"doesExist,omitempty"`
	Error    string `json:"error,omitempty"`
	Message  string `json:"message,omitempty"`
	Replaced *bool  `json:"replaced,omitempty"`
	Value    string `json:"value,omitempty"`
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
