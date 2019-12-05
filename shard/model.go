package shard

import (
	"github.com/colbyleiske/cse138_assignment2/vectorclock"
)

type Shard struct {
	ID          string // kept as string for easier consistent hashing sorting
	Nodes       []string
	VectorClock *vectorclock.VectorClock
}
