package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ainghazal/kurma"
)

var (
	port = ":8080"
)

func decoyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello stranger"))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello friend"))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: server secret|nonce")
		os.Exit(1)
	}
	mode := os.Args[1]
	st := -1
	switch mode {
	case "secret":
		st = kurma.StaticSecret
	case "nonce":
		st = kurma.ValidNonce
	default:
		fmt.Println("unknown mode", mode)
		os.Exit(1)
	}
	log.Println("Starting server at port", port)

	dcr := kurma.Discriminator{
		Strategy:  st,
		Decoy:     decoyHandler,
		Protected: protectedHandler,
	}

	http.HandleFunc("/", kurma.WithDecoyHandler(dcr))

	err := http.ListenAndServeTLS(port, "./testdata/server.crt", "./testdata/server.key", nil)
	if err != nil {
		panic(err)
	}
}
