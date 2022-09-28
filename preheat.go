package kurma

import (
	"crypto/sha512"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

const (
	tokenLen = 142
)

func PreHeat(authURI, targetURI string) net.Conn {
	nonce := fetchNonce(authURI)
	fmt.Printf("nonce: %s\n", nonce)
	return preHeater(nonce, targetURI)
}

func fetchNonce(uri string) string {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", uri, nil)
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Println("[!] No TOKEN set!!!")
		ts := time.Now().Unix()
		hsh := fmt.Sprintf("%x", sha512.Sum512([]byte(string(ts))))
		token = (hsh + hsh)[:tokenLen]
	}
	req.Header.Set(secretHeader, token)
	resp, err := client.Do(req)
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
