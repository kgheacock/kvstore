package hasher

import "testing"

var ring Ring

func init() {
	//setup ring
	//add nodes
	//add keys
}

func TestGetServerByKey(t *testing.T) {
	ring := NewRing()
	tt := []struct {
		key           string
		expectedIP    string
		expectedError error
	}{
		{key: "foo", expectedIP: "127.0.0.1", expectedError: nil},
		{key: "bar", expectedIP: "", expectedError: KeyNotFound},
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
}
