package hasher

import (
	"testing"
)

var ring Ring

func TestBasicCorrectMapping(t *testing.T) {
	ring := NewRing()
	servers := []string{"A", "B", "C"}
	for item := range servers {
		ring.AddServer(servers[item])
	}

	tt := []struct {
		key           string
		expectedIP    string
		expectedError error
	}{
		{key: "Chris", expectedIP: "C", expectedError: nil},
		{key: "Brandon", expectedIP: "C", expectedError: nil},
		{key: "Colby", expectedIP: "B", expectedError: nil},
		{key: "Keith", expectedIP: "B", expectedError: nil},
		{key: "Alvaro", expectedIP: "A", expectedError: nil},
		{key: "Mackey", expectedIP: "A", expectedError: nil},
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

func TestGetServerByKey(t *testing.T) {
	ring := NewRing()
	servers := []string{"A", "B", "C"}
	for item := range servers {
		ring.AddServer(servers[item])
	}

	tt := []struct {
		key           string
		expectedIP    string
		expectedError error
	}{
		{key: "Tantalo", expectedIP: "C", expectedError: nil},
		{key: "Alvaro", expectedIP: "A", expectedError: nil},
		{key: "DoesNotExist", expectedIP: "C", expectedError: nil},
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
