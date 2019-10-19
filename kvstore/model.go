package kvstore

//Store ...
type Store struct {
	dal DataAccessLayer
}

//Anonymous struct used for the ResponseMessage struct
//This is neccessary for use with omitempty in json object
exists := struct {
	Exists bool
}

//Anonymous struct used for the ResponseMessage struct
//This is neccessary for use with omitempty in json object
replaced := struct {
	Replaced bool
}

//Error and Success message
type ResponseMessage struct {
	Exists  *exists  `json:"doesExist,omitempty"`
	Error   string  `json:"error,omitempty"`
	Message string  `json:"message,omitempty"`
	Replaced *replaced `json:"replaced,omitempty"`
	Value string    `json:"value,omitempty"`
}

//DataAccessLayer interface
type DataAccessLayer interface {
	Delete(key string) error

	Put(key string, value string) (string, error)

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
