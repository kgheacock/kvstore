package main

import (
	"github.com/colbyleiske/cse138_assignment2/kvstore"

)

func main() {
	dal := kvstore.KVDAL{}
	s := kvstore.NewStore(&dal)

	s.DAL().Delete()
}