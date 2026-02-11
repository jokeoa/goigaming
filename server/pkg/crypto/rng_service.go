package crypto

import (
	"github.com/jokeoa/goigaming/internal/core/domain"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GenerateServerSeed() (string, error) {
	return GenerateServerSeed()
}

func (s *Service) HashSeed(serverSeed string) string {
	return HashSeed(serverSeed)
}

func (s *Service) ShuffleDeck(serverSeed, clientSeed string, nonce int) []domain.Card {
	return ShuffleDeck(serverSeed, clientSeed, nonce)
}

func (s *Service) VerifySeed(serverSeed, hash string) bool {
	return VerifySeed(serverSeed, hash)
}
