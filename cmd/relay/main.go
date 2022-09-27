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

func main() {
	log.Println("Starting server at port", port)

	dcr := kurma.Discriminator{
		Strategy:  kurma.ValidNonce,
		Decoy:     decoyHandler,
		Protected: kurma.HijackAndProxyHandler("127.0.0.1:4430"),
	}

	http.HandleFunc("/", kurma.WithDecoyHandler(dcr))

	err := http.ListenAndServeTLS(port, "./testdata/server.crt", "./testdata/server.key", nil)
	if err != nil {
		panic(err)
	}
}
