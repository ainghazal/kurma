package main

import (
	"fmt"
	"os"
)

func main() {
	key := loadKey("key")
	token := os.Getenv("TOKEN")
	valid := VerifyToken(token, key)
	fmt.Println("signature verified:", valid)
}
