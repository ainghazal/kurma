package main

import (
	"crypto/ecdsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
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

func loadKey(file string) *ecdsa.PrivateKey {
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

func VerifyToken(token string, key *ecdsa.PrivateKey) bool {
	hash := sha1.Sum([]byte(secret))
	sig, err := hex.DecodeString(token)
	if err != nil {
		panic(err)
	}
	return ecdsa.VerifyASN1(&key.PublicKey, hash[:], sig)
}
