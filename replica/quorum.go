package replica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/colbyleiske/cse138_assignment2/config"
	"github.com/colbyleiske/cse138_assignment2/kvstore"
)

func PutQuorum(key string, value string) error {
	thisQuorum := config.Config.ThisQuorom
	replFactor := config.Config.ReplFactor
	ch := make(chan kvstore.PutResponse)
	ips := config.Config.Quoroms[thisQuorum]

	for _, ip := range ips {
		// read to all quorums
		if thisQuorum != ip {
			go func() {
				body := fmt.Sprintf("{\"value\":%s}", value)
				resp, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s/replication/%s", ip, key), strings.NewReader(body))
				decoder := json.NewDecoder(resp.Body)
				var data kvstore.PutResponse
				decoder.Decode(&data)
				ch <- data
			}()
		}
	}

	for len(ch) > replFactor/2 {
		// wait for N/2 > acks then check for most recent one
	}

	for len(ch) > 0 {
		var data kvstore.PutResponse
		data = <-ch
		fmt.Println(data.Value)
	}
	// Get Request with the highest vector clock

	return nil
}

func GetQuorum(key string) error {
	thisQuorum := config.Config.ThisQuorom
	replFactor := config.Config.ReplFactor
	ch := make(chan kvstore.GetResponse)
	ips := config.Config.Quoroms[thisQuorum]

	for _, ip := range ips {
		// read to all quorums
		if thisQuorum != ip {
			go func() {
				resp, err := http.Get(fmt.Sprintf("http://%s/replication/%s", ip, key))
				decoder := json.NewDecoder(resp.Body)
				var data kvstore.GetResponse
				decoder.Decode(&data)
				ch <- data
			}()
		}
	}

	for len(ch) > replFactor/2 {
		// wait for N/2 > acks then check for most recent one
	}

	for len(ch) > 0 {
		var data kvstore.GetResponse
		data = <-ch
		fmt.Println(data.Value)
	}

	// Get Request with the highest vector clock

	return nil
}

func DeleteQuorum(key string) error {
	thisQuorum := config.Config.ThisQuorom
	replFactor := config.Config.ReplFactor
	ch := make(chan kvstore.DeleteResponse)
	ips := config.Config.Quoroms[thisQuorum]

	for _, ip := range ips {
		// read to all quorums
		if thisQuorum != ip {
			go func() {
				resp, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s/replication/%s", ip, key), nil)
				decoder := json.NewDecoder(resp.Body)
				var data kvstore.DeleteResponse
				decoder.Decode(&data)
				ch <- data
			}()
		}
	}

	for len(ch) > replFactor/2 {
		// wait for N/2 > acks then check for most recent one
	}

	for len(ch) > 0 {
		var data kvstore.DeleteResponse
		data = <-ch
		fmt.Println(data.Exists)
	}
	// Get Request with the highest vector clock

	return nil
}
