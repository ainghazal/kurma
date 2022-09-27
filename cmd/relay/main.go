package main

import (
	"log"
	"net/http"

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
	log.Println("Starting server at port", port)

	dcr := kurma.Discriminator{
		Strategy:  kurma.ValidNonce,
		Decoy:     decoyHandler,
		Protected: protectedHandler,
	}

	http.HandleFunc("/", kurma.DecoyHandler(dcr))

	err := http.ListenAndServeTLS(port, "./testdata/server.crt", "./testdata/server.key", nil)
	if err != nil {
		panic(err)
	}
}
