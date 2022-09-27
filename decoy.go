package kurma

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	magicCookie  = "letmein"
	secretHeader = "X-CSRF-Token"
	authBaseURL  = "http://localhost:3001/nonce/"
)

type discriminatorFn = func(string) bool

func isRequestWithMagicCookie(val string) bool {
	return bytes.Equal([]byte(strings.TrimSpace(val)), []byte(magicCookie))
}

func isRequestSigned(val string) bool {
	// TODO implement
	return false
}

func isValidNonce(val string) bool {
	checkURI := authBaseURL + val + "/"
	resp, err := http.Get(checkURI)
	if err == nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil && strings.TrimSpace(string(body)) == "true" {
			return true
		}
	}
	return false
}

const (
	StaticSecret = iota
	EcdsaSignature
	ValidNonce
)

type Discriminator struct {
	Strategy  int
	Decoy     http.HandlerFunc
	Protected http.HandlerFunc
}

func DecoyHandler(d Discriminator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := r.Header.Get(secretHeader)
		log.Println(secretHeader+":", cookie)

		var evalFn discriminatorFn
		switch d.Strategy {
		case StaticSecret:
			evalFn = isRequestWithMagicCookie
		case EcdsaSignature:
			evalFn = isRequestSigned
		case ValidNonce:
			evalFn = isValidNonce
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("bad gateway"))
			return
		}

		if !evalFn(cookie) {
			d.Decoy(w, r)
			return
		} else {
			d.Protected(w, r)
			return
			/*
			 conn, _, err := res.(http.Hijacker).Hijack()
			 if err != nil {
			 	panic(err)
			 }
			 conn.Write([]byte{})
			 fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n")
			*/

			/*
			 buffer := make([]byte, 1024)
			 fmt.Println("Server : Enter routine")
			 for {
			 	time.Sleep(1 * time.Second)
			 	fmt.Println("Server : I send")
			 	_, err = conn.Write([]byte("Hijack server"))
			 	if err != nil {
			 		panic(err)
			 	}
			 	fmt.Println("Server : I'm receiving")
			 	n, err := conn.Read(buffer)
			 	if err != nil {
			 		panic(err)
			 	}
			 	fmt.Printf("Server : %d bytes from client : %s\n", n, string(buffer))
			 }
			*/
		}
	}
}
