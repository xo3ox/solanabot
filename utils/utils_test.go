package utils

import (
	"testing"
)

func TestGetBalance(t *testing.T) {
	if err := GetBalance(); err != nil {
		t.Errorf("GetBalance() error = %v", err)
	}
}

func TestNewRandomPrivateKey(t *testing.T) {
	// Generate a new random private key
	NewRandomPrivateKey()
}
