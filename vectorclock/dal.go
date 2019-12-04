package vectorclock

import (
	"fmt"

	"github.com/colbyleiske/cse138_assignment2/config"
)

//VectorClock contains VC a map, and the current servers VC value
type VectorClock struct {
	VC map[string]int
	ip string
}

//Len returns len() of map of a VectorClock
func (vc *VectorClock) Len() int { return len(vc.VC) }

//NewVectorClock creates VectorClock object
func NewVectorClock() *VectorClock {
	m := make(map[string]int)
	server := config.Config.Address
	m[server] = 0
	return &VectorClock{VC: m, ip: server}
}

//IncrementVC is called on every succesful Put and Get
func (vc *VectorClock) IncrementVC() {
	server := config.Config.Address
	vc.VC[server]++
}

//CurrentState returns value of own clock
func (vc *VectorClock) CurrentState() int {
	return vc.VC[config.Config.Address]
}

//IP returns the IP of the server with that VC
func (vc *VectorClock) IP() string {
	return vc.ip
}

//UpdateVC updates vc's values by taking piecewise max. Only changes vc.
func (vc *VectorClock) UpdateVC(vc2 *VectorClock) {
	serverList := config.Config.Quoroms[config.Config.ThisQuorom]
	for _, server := range serverList {
		if vc.VC[server] < vc2.VC[server] {
			vc.VC[server] = vc2.VC[server]
		}
	}
}

//MaxClock returns IP of max-clock server from slice
func MaxClock(vclist []*VectorClock) string {
	max := 0
	//If all clocks equal, returns IP of first in slice
	maxserverIP := vclist[0].IP()
	for _, item := range vclist {
		current := 0
		vcmap := item.VC
		for _, v := range vcmap {
			current += v
		}
		if max < current {
			max = current
			maxserverIP = item.IP()
		}
	}
	return maxserverIP
}

//Print is a debugging function
func (vc *VectorClock) Print() {
	serverList := config.Config.Quoroms[config.Config.ThisQuorom]
	fmt.Print("[")
	for _, server := range serverList {
		fmt.Printf("%v, ", vc.VC[server])
	}
	fmt.Print(config.Config.Address)
	fmt.Println("]")
}
