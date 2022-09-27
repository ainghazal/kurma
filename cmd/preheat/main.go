package main

import (
	"log"

	"github.com/ainghazal/kurma"
)

var (
	authURI   = "http://localhost:3000"
	targetURI = "localhost:8080"
)

func main() {
	conn := kurma.PreHeat(authURI, targetURI)
	log.Println("got conn", conn)
}

// dialer:
//./gost -L=tcp://127.0.0.1:2222/37.218.244.248:1194 -F='obfs4://172.82.146.246:443/?cert=avWz75FFdYxm4wzrz9mJDJRK48l4Ip7KOLilMWHc%2BwvfRg%2FmBsDb1xFeZiPwxWWypjeDcQ&iat-mode=2'
