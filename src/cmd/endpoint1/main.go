package main

import (
	"log"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/router"
)

var version = "0"

func main() {
	log.Printf("DPV version %s", version)
	server := router.NewServer("config.yml", false)
	dpv.ConfigInstance.Settings.Version = version
	log.Fatal(server.ListenAndServe())
}
