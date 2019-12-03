package vectorclock

import (
	"fmt"

	"github.com/colbyleiske/cse138_assignment2/config"
)

//VectorClock contains VC a map, and the current servers VC value
type VectorClock struct {
	VC map[string]int
}

//Less compares two vectorclocks value at the same server key, "if v1 less than v2"
//func (vc *VectorClock) Less(vc2 *VectorClock) bool {
//	return vc.VC[config.Config.Address] < vc2.VC[config.Config.Address]
//}

//Len returns len() of map of a VectorClock
func (vc *VectorClock) Len() int { return len(vc.VC) }

//NewVectorClock creates VectorClock object
func NewVectorClock() *VectorClock {
	m := make(map[string]int)
	server := config.Config.Address
	m[server] = 0
	return &VectorClock{VC: m}
}

//IncrementVC is called on every succesful Put and Get
func (vc *VectorClock) IncrementVC() {
	server := config.Config.Address
	vc.VC[server]++
}

//CurrentState returns int currentState from VC struct
func (vc *VectorClock) CurrentState() int {
	return vc.VC[config.Config.Address]
}

//UpdateVC matches two VC's by taking piecewise max of both
func (vc *VectorClock) UpdateVC(vc2 *VectorClock) {
	serverList := config.Config.Servers
	for _, server := range serverList {
		if vc.VC[server] < vc2.VC[server] {
			vc.VC[server] = vc2.VC[server]
		}
	}
}

//Print is a debugging function
func (vc *VectorClock) Print() {
	serverList := config.Config.Servers
	fmt.Print("[")
	for _, server := range serverList {
		fmt.Printf("%v, ", vc.VC[server])
	}
	fmt.Print(config.Config.Address)
	fmt.Println("]")
}
