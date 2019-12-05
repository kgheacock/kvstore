package shard

import "github.com/colbyleiske/cse138_assignment2/vectorclock"

type Shard struct {
	ID          int
	Nodes       []string
	VectorClock *vectorclock.VectorClock
}
