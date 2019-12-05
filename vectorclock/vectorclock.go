package vectorclock

//VectorClock contains VC a map, and the current servers VC value
type VectorClock struct {
	Clocks map[string]int
	ip     string
}

//Len returns len() of map of a VectorClock
func (vc *VectorClock) Len() int { return len(vc.Clocks) }

//NewVectorClock creates VectorClock object
func NewVectorClock(ips []string, localIP string) *VectorClock {
	vc := &VectorClock{ip: localIP}
	vc.ResetVC(ips)

	return vc
}

//IncrementVC is called on every succesful Put and Get
func (vc *VectorClock) IncrementNodeClock() {
	vc.Clocks[vc.ip]++
}

//CurrentState returns value of own clock
func (vc *VectorClock) CurrentServerClock() int {
	return vc.Clocks[vc.ip]
}

//ResetVC resets a VC map to 0
func (vc *VectorClock) ResetVC(serverList []string) {
	vc.Clocks = make(map[string]int)
	for _, server := range serverList {
		vc.Clocks[server] = 0
	}
}

//IP returns the IP of the server with that VC
func (vc *VectorClock) IP() string {
	return vc.ip
}

//UpdateVC updates vc's values by taking piecewise max. Only changes vc.
func (vc *VectorClock) UpdateClocks(incVC *VectorClock) {
	//Can assume incoming vector clock has same servers as us.
	for nodeIP := range vc.Clocks {
		if vc.Clocks[nodeIP] < incVC.Clocks[nodeIP] {
			vc.Clocks[nodeIP] = incVC.Clocks[nodeIP]
		}
	}
}

func (vc *VectorClock) HappenedBefore (incVC *VectorClock) bool {
	for ip := range vc.Clocks {
		if vc.Clocks[ip] > incVC.Clocks[ip] {
			return false
		}
	}

	return true
}

//ReceiveEvent is called whenever an event is delivered and we need to tick + update our clock.
func (vc *VectorClock) ReceiveEvent(incVC *VectorClock) {
	vc.IncrementNodeClock()
	vc.UpdateClocks(incVC)
}

// //Print is a debugging function
// func (vc *VectorClock) Print() {
// 	serverList := config.Config.Quoroms[config.Config.ThisQuorom]
// 	fmt.Print("[")
// 	for _, server := range serverList {
// 		fmt.Printf("%v, ", vc.VC[server])
// 	}
// 	fmt.Print(config.Config.Address)
// 	fmt.Println("]")
// }
