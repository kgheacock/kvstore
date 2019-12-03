package vectorclock

type VectorClock struct {
	VC      map[string]int
	MyValue int
}

//NewVectorClock creates VectorClock object
func NewVectorClock() *VectorClock {
	return &VectorClock{VC: make(map[string]int), MyValue: 0}
}

//IncrementVC is called on every succesful Put and Get
func (vc *VectorClock) IncrementVC() {

}

//CompareVC compares two vector clocks
func (vc *VectorClock) CompareVC() {

}

//Updates single value in VectorClock map
func (vc *VectorClock) updatentry() {

}
