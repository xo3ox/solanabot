package utils

import (
	"fmt"
	"testing"
)

func TestNewRandomPrivateKey(t *testing.T) {
	// Generate a new random private key
	NewRandomPrivateKey()

}

func TestGetKeypair(t *testing.T) {
	privateKey, publicKey := GetKeypair("../.env")
	fmt.Println("Private Key:", privateKey)
	fmt.Println("Private Key:", []byte(privateKey))
	fmt.Println("Public Key:", publicKey)
}
