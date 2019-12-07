package shard

type Shard struct {
	ID    string // kept as string for easier consistent hashing sorting
	Nodes []string
}
