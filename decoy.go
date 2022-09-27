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
	ValidNonce
)

type Discriminator struct {
	Strategy  int
	Decoy     http.HandlerFunc
	Protected http.HandlerFunc
}

func WithDecoyHandler(d Discriminator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := r.Header.Get(secretHeader)
		log.Println(secretHeader+":", cookie)

		var evalFn discriminatorFn
		switch d.Strategy {
		case StaticSecret:
			evalFn = isRequestWithMagicCookie
		case ValidNonce:
			evalFn = isValidNonce
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("bad gateway"))
			return
		}
		if !evalFn(cookie) {
			log.Println("Handing decoy...")
			d.Decoy(w, r)
			return
		} else {
			log.Println("A friend knocked...")
			d.Protected(w, r)
			return
		}
	}
}
