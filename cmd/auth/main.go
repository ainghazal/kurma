package main

import (
	"fmt"
	"net/http"

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

func defaultHandler(c echo.Context) error {
	request := new(APIRequest)
	binder := &echo.DefaultBinder{}
	binder.BindHeaders(c, request)
	fmt.Printf("%+v\n", request)

	// just testing...
	// TODO pass a true condition to return this handler in a closure

	nonce := nonces.New()
	c.Response().Header().Set(hdr, nonce)

	return c.String(http.StatusOK, fmt.Sprintf("%s\n", c.RealIP()))
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
	e := echo.New()
	e.IPExtractor = echo.ExtractIPDirect()
	e.GET("/", defaultHandler)
	go e.Start(externalAuthURI)

	el := echo.New()
	el.GET("/nonce/:id/", validNonceHandler)
	el.Start(localQueryURI)
}
