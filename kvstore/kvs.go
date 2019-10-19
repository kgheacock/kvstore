package kvstore

import (
	"net/http"

	"github.com/gorilla/mux"
)

//http endpoints go here
func (s *Store) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		w.Write([]byte("something something bad request"))
	}
	s.DAL().Delete(key)
}
func (s *Store) AddHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *Store) GetHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Store) UpdateHandler(w http.ResponseWriter, r *http.Request) {
}
