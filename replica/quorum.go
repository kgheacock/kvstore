package replica

type Quorum struct {
	id       [16]byte
	ips      []string
	quorumid string
}

// make seed all the ips of the view

func (q *Quorum) PutQuorum(key string, value string) error {
	acks := make(chan string, len(q.ips)-1)

	for i := 0; i < len(q.ips); i++ {
		// write to all quorums
		if q.quorumid != q.ips[i] {
			acks <- "GET" // http.Get(fmt.Sprintf("%s/replication/%s", q.ips[i], key, url.Values{"value": {value}}))
		}

	}

	for len(acks) <= len(q.ips)/2 {
		return
	}
}

func (q *Quorum) GetQuorum(key string) error {
	for i := 0; i < len(q.ips); i++ {
		// read to all quorums
	}

	// wait for N/2 > acks then check for most recent one

	return nil
}

func (q *Quorum) DeleteQuorum(key string) error {
	for i := 0; i < len(q.ips); i++ {
		// read to all quorums
	}

	// wait for N/2 > acks then check for most recent one

	return nil
}
