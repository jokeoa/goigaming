package game

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/jokeoa/goigaming/internal/core/domain"
)

type SimpleRNGService struct{}

func (s *SimpleRNGService) GenerateServerSeed() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("SimpleRNGService.GenerateServerSeed: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func (s *SimpleRNGService) HashSeed(serverSeed string) string {
	hash := sha256.Sum256([]byte(serverSeed))
	return hex.EncodeToString(hash[:])
}

func (s *SimpleRNGService) ShuffleDeck(serverSeed, clientSeed string, nonce int) []domain.Card {
	deck := domain.FullDeck()

	combined := fmt.Sprintf("%s:%s:%d", serverSeed, clientSeed, nonce)
	seed := sha256.Sum256([]byte(combined))

	for i := len(deck) - 1; i > 0; i-- {
		max := big.NewInt(int64(i + 1))
		hashInput := fmt.Sprintf("%x:%d", seed, i)
		h := sha256.Sum256([]byte(hashInput))
		n := new(big.Int).SetBytes(h[:])
		j := int(new(big.Int).Mod(n, max).Int64())

		deck[i], deck[j] = deck[j], deck[i]
	}

	return deck
}

func (s *SimpleRNGService) VerifySeed(serverSeed, hash string) bool {
	return s.HashSeed(serverSeed) == hash
}
