package utils

import (
	"testing"
)

func TestNewRandomPrivateKey(t *testing.T) {
	// Generate a new random private key
	NewRandomPrivateKey()

}

func TestGetKeypair(t *testing.T) {
	_ = GetKeypair("../.env")
}
