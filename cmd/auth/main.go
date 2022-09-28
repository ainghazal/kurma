package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ainghazal/kurma"
	"github.com/labstack/echo/v4"
)

var (
	externalAuthURI = ":3000"
	localQueryURI   = "localhost:3001"
	hdr             = "X-Csrf-Token"
)

type APIRequest struct {
	Token string `header:"X-CSRF-TOKEN"`
}

var nonces = &kurma.NonceJar{}

func allowAll(string) bool {
	return true
}

func denyAll(string) bool {
	return false
}

func allowSignedBy(key *ecdsa.PrivateKey) func(string) bool {
	return func(token string) bool {
		return kurma.VerifyToken(token, key)
	}
}

func defaultHandler(debug bool, key *ecdsa.PrivateKey) func(echo.Context) error {
	return func(c echo.Context) error {
		c.Response().Header().Set("Server", "nginx/1.14.1")
		request := new(APIRequest)
		binder := &echo.DefaultBinder{}
		binder.BindHeaders(c, request)

		authFn := denyAll
		switch debug {
		case true:
			log.Println("[+] debug mode")
			authFn = allowAll
		default:
			log.Println("[+] allow-signed mode")
			authFn = allowSignedBy(key)
		}

		var nonce string
		switch authFn(request.Token) {
		case true:
			if request.Token != "" {
				log.Println("got valid token:", request.Token)
			}
			nonce = nonces.New()
		default:
			nonce = kurma.FakeNonce()
		}
		c.Response().Header().Set(hdr, nonce)
		return c.String(http.StatusOK, fmt.Sprintf("%s\n", c.RealIP()))
	}
}

func validNonceHandler(c echo.Context) error {
	switch nonces.IsValid(c.Param("id")) {
	case true:
		return c.String(http.StatusOK, "true\n")
	default:
		return c.String(http.StatusOK, "false\n")
	}
}

func main() {
	// TODO configure external and api ports
	debug := os.Getenv("DEBUG") == "1"
	key := kurma.LoadKey("testdata/key")

	e := echo.New()
	e.IPExtractor = echo.ExtractIPDirect()
	e.GET("/", defaultHandler(debug, key))
	go e.Start(externalAuthURI)

	el := echo.New()
	el.GET("/nonce/:id/", validNonceHandler)
	el.Start(localQueryURI)
}
