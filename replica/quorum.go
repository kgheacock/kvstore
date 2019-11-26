package replica

type Quorum struct {
	id  [16]byte
	ips []string
}

// make seed all the ips of the view

func (q *Quorum) WriteQuorum(key string) error {
	for i := 0; i < len(q.ips); i++ {
		// write to all quorums
	}

	// wait for N/2 > acks then exits

	return nil
}

func (q *Quorum) ReadQuorum(key string) error {
	for i := 0; i < len(q.ips); i++ {
		// read to all quorums
	}

	// wait for N/2 > acks then check for most recent one

	return nil
}
