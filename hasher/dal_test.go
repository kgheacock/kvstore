package hasher

import (
	"testing"
)

var ring Ring

func init() {

}

func TestGetServerByKey(t *testing.T) {
	keyset1 := []string{"Chris", "Brandon", "Colby", "Keith", "Alvaro", "Mackey", "Space-Boss", "Dimitris", "Julig", "Pham", "Betz", "Long", "Qian", "Tantalo"}
	servers := []string{"A", "B", "C"}
	//setup ring
	//add nodes
	//add keys

	ring := NewRing()
	for item := range keyset1 {
		ring.AddKey(keyset1[item])
	}
	for item := range servers {
		ring.AddServer(servers[item])
	}
	ring.ReShard()

	tt := []struct {
		key           string
		expectedIP    string
		expectedError error
	}{
		{key: "Tantalo", expectedIP: "C", expectedError: nil},
		{key: "Alvaro", expectedIP: "A", expectedError: nil},
		{key: "DoesNotExist", expectedIP: "", expectedError: KeyNotFound},
	}

	for _, tc := range tt {
		ip, err := ring.GetServerByKey(tc.key)
		if ip != tc.expectedIP {
			t.Errorf("Expected IP %s , got %s", tc.expectedIP, ip)
		}
		if err != tc.expectedError {
			t.Errorf("Expected error %s , got %s", tc.expectedError, err)

		}
	}

	allKeys := GetServersAndKeys()
	for i := 0; i < len(allKeys); i++ {
		t.Logf(allKeys[i].ServerName)
	}
}
