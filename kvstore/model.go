package kvstore

type Store struct {
	dal DataAccessLayer
}

type DataAccessLayer interface {
	Delete(key string) error
	Get(key string) (string, error)
	Put(key string, value string) (int, error)
}

func NewStore(dal DataAccessLayer) *Store {
	return &Store{dal: dal}
}

func (s *Store) DAL() DataAccessLayer {
	return s.dal
}

//Holds incoming PUT request body
type Data struct {
	Value string `json:"value"`
}

type ResponseMessage struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Value   string `json:"value,omitempty"`
}

type DeleteResponse struct {
	ResponseMessage
	Exists bool `json:"doesExist"`
}

type PutResponse struct {
	ResponseMessage
	Replaced bool `json:"replaced"`
}

type GetResponse struct {
	ResponseMessage
	Exists bool `json:"doesExist"`
}
