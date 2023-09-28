package main

import (
	"log"
	"pkv/api/src/router"
)

func main() {
	server := router.NewServer("config.yml")
	log.Fatal(server.ListenAndServe())
}
