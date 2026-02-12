package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GenerateServerSeed() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("crypto.GenerateServerSeed: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func HashSeed(serverSeed string) string {
	h := sha256.Sum256([]byte(serverSeed))
	return hex.EncodeToString(h[:])
}

func VerifySeed(serverSeed, hash string) bool {
	return HashSeed(serverSeed) == hash
}
