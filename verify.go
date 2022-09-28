package kurma

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const secret = "trust is never trust enough"

func encodeKey(privateKey *ecdsa.PrivateKey) string {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return string(pemEncoded)
}

func decodeKey(pemEncoded string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)
	return privateKey
}

func LoadKey(file string) *ecdsa.PrivateKey {
	keyFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	keyBytes, err := ioutil.ReadAll(keyFile)
	if err != nil {
		log.Fatal(err)
	}
	return decodeKey(string(keyBytes))

}

func genKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(encodeKey(privateKey))
	return privateKey
}

func LoadOrGenKey(keyfile string) *ecdsa.PrivateKey {
	_, err := os.Stat(keyfile)
	if errors.Is(err, os.ErrNotExist) {
		return genKey()
	}
	keyFile, err := os.Open(keyfile)
	if err != nil {
		log.Fatal(err)
	}
	keyBytes, err := ioutil.ReadAll(keyFile)
	if err != nil {
		log.Fatal(err)
	}
	return decodeKey(string(keyBytes))

}

func VerifyToken(token string, key *ecdsa.PrivateKey) bool {
	hash := sha1.Sum([]byte(secret))
	sig, err := hex.DecodeString(token)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return ecdsa.VerifyASN1(&key.PublicKey, hash[:], sig)
}

func GenerateTokens(num int, key *ecdsa.PrivateKey) {
	hash := sha1.Sum([]byte(secret))

	for i := 0; i < num; i++ {
		sig, err := ecdsa.SignASN1(rand.Reader, key, hash[:])
		if err != nil {
			panic(err)
		}
		fmt.Printf("%x\n", sig)
	}
}
