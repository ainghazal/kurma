package main

import (
	"fmt"
	"os"

	"github.com/ainghazal/kurma"
)

func verify() {
	key := kurma.LoadKey("testdata/key")
	token := os.Getenv("TOKEN")
	valid := kurma.VerifyToken(token, key)
	fmt.Println("signature verified:", valid)
}

func sign(keyfile string) {
	key := kurma.LoadOrGenKey(keyfile)
	kurma.GenerateTokens(20, key)
}

func usage() {
	fmt.Printf("Usage: %s verify|sign [keyfile]\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	cmd := os.Args[1]
	switch cmd {
	case "verify":
		verify()
		return
	case "sign":
		if len(os.Args) < 3 {
			usage()
		}
		key := os.Args[2]
		sign(key)
		return
	default:
		usage()
	}
}
