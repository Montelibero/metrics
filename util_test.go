package metrics

import (
	"testing"

	"github.com/stellar/go/protocols/horizon"
)

func TestIsDataEqual(t *testing.T) {
	type testCase struct {
		name     string
		acc      horizon.Account
		key      string
		value    string
		expected bool
	}

	testCases := []testCase{
		{
			name: "Data matches",
			acc: horizon.Account{
				Data: map[string]string{
					"testKey": base64Encode("testValue"),
				},
			},
			key:      "testKey",
			value:    "testvalue",
			expected: true,
		},
		{
			name: "Data does not match",
			acc: horizon.Account{
				Data: map[string]string{
					"testKey": base64Encode("differentValue"),
				},
			},
			key:      "testKey",
			value:    "testvalue",
			expected: false,
		},
		{
			name: "Key does not exist",
			acc: horizon.Account{
				Data: map[string]string{},
			},
			key:      "nonExistentKey",
			value:    "testvalue",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isDataEqual(tc.acc, tc.key, tc.value)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestIsDataExist(t *testing.T) {
	type testCase struct {
		name     string
		acc      horizon.Account
		key      string
		expected bool
	}

	testCases := []testCase{
		{
			name: "Data exists",
			acc: horizon.Account{
				Data: map[string]string{
					"testKey": base64Encode("testValue"),
				},
			},
			key:      "testKey",
			expected: true,
		},
		{
			name: "Data does not exist",
			acc: horizon.Account{
				Data: map[string]string{},
			},
			key:      "testKey",
			expected: false,
		},
		{
			name: "Key does not exist",
			acc: horizon.Account{
				Data: map[string]string{
					"anotherKey": base64Encode("anotherValue"),
				},
			},
			key:      "testKey",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isDataExist(tc.acc, tc.key)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
