package kurma

import (
	"bytes"
	"crypto/rand"
	"log"
	"math/big"
	"sync"
	"time"
)

var (
	nonceLen             = 8
	tokenLifetimeSeconds = 30 // TODO make this configurable (tune it in production)
)

func FakeNonce() string {
	str, _ := generateRandomString(nonceLen)
	return str
}

type NonceJar struct {
	tokens []string
	mu     sync.Mutex
}

func (n *NonceJar) New() string {
	n.mu.Lock()
	defer n.mu.Unlock()
	nonce, _ := generateRandomString(nonceLen)
	log.Println("adding nonce", nonce)
	n.tokens = append(n.tokens, nonce)
	log.Printf("got %d tokens\n", len(n.tokens))
	go n.expire(nonce)
	return nonce
}

func (n *NonceJar) expire(token string) {
	time.AfterFunc(time.Second*time.Duration(tokenLifetimeSeconds), func() {
		n.mu.Lock()
		defer n.mu.Unlock()
		for i, nc := range n.tokens[:] {
			if bytes.Equal([]byte(token), []byte(nc)) {
				// we delete the expired token from the jar
				n.tokens = append(n.tokens[:i], n.tokens[i+1:]...)
				log.Printf("expired nonce %s\n", token)
				return
			}
		}
	})
}

func (n *NonceJar) IsValid(token string) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	for i, nc := range n.tokens[:] {
		if bytes.Equal([]byte(token), []byte(nc)) {
			// we delete the used token from the jar
			n.tokens = append(n.tokens[:i], n.tokens[i+1:]...)
			log.Printf("consumed nonce %s\n", token)
			return true
		}
	}
	return false
}

func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, 0)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret = append(ret, letters[num.Int64()])
	}
	return string(ret), nil
}
