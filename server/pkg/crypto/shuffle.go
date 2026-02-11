package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/jokeoa/goigaming/internal/core/domain"
)

func ShuffleDeck(serverSeed, clientSeed string, nonce int) []domain.Card {
	deck := domain.FullDeck()
	gen := newEntropyGenerator(serverSeed, clientSeed, nonce)

	for i := len(deck) - 1; i > 0; i-- {
		j := gen.unbiasedRandom(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}

	return deck
}

type entropyGenerator struct {
	serverSeed string
	clientSeed string
	nonce      int
	round      int
	buffer     []uint32
	bufferIdx  int
}

func newEntropyGenerator(serverSeed, clientSeed string, nonce int) *entropyGenerator {
	return &entropyGenerator{
		serverSeed: serverSeed,
		clientSeed: clientSeed,
		nonce:      nonce,
	}
}

func (g *entropyGenerator) next() uint32 {
	for g.bufferIdx >= len(g.buffer) {
		g.buffer = g.generateBlock()
		g.bufferIdx = 0
	}
	val := g.buffer[g.bufferIdx]
	g.bufferIdx++
	return val
}

func (g *entropyGenerator) generateBlock() []uint32 {
	message := fmt.Sprintf("%s:%d:%d", g.clientSeed, g.nonce, g.round)
	h := hmac.New(sha256.New, []byte(g.serverSeed))
	h.Write([]byte(message))
	hash := h.Sum(nil)
	g.round++

	results := make([]uint32, 0, len(hash)/4)
	for i := 0; i+4 <= len(hash); i += 4 {
		results = append(results, binary.BigEndian.Uint32(hash[i:i+4]))
	}
	return results
}

func (g *entropyGenerator) unbiasedRandom(max int) int {
	m := uint32(max)
	threshold := (0xFFFFFFFF - m + 1) % m

	for {
		val := g.next()
		if val >= threshold {
			return int(val % m)
		}
	}
}

func hexToUint32(hexStr string) (uint32, error) {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b), nil
}
