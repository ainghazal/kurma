package kurma

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

func PreHeat(authURI, targetURI string) net.Conn {
	nonce := fetchNonce(authURI)
	fmt.Printf("nonce: %s\n", nonce)
	return preHeater(nonce, targetURI)
}

func fetchNonce(uri string) string {
	// TODO set secret here for auth endpoint.
	resp, err := http.Get(uri)
	if err != nil {
		panic(err)
	}
	return resp.Header.Get(secretHeader)
}

func preHeater(nonce, uri string) net.Conn {
	req, err := http.NewRequest("GET", "/", nil)
	req.Header.Set(secretHeader, nonce)
	if err != nil {
		panic(err)
	}

	log.Println("Dialing:", uri)
	dial, err := net.Dial("tcp", uri)
	if err != nil {
		panic(err)
	}
	log.Println("Client: create TLS connection")

	// TODO can get the server name from uri, or as an optional parameter.
	tls_conn := tls.Client(dial, &tls.Config{ServerName: "localhost"})
	conn := httputil.NewClientConn(tls_conn, nil)

	_, err = conn.Do(req)
	if err != httputil.ErrPersistEOF && err != nil {
		panic(err)
	}

	log.Println("Client: hijacking https connection")
	connection, _ := conn.Hijack()
	return connection
}
