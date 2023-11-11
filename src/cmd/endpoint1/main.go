package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/router"
)

var version = "0"

func main() {
	log.Printf("DPV version %s", version)
	server := router.NewServer("config.yml", false)
	dpv.ConfigInstance.Settings.Version = version
	socketPath := os.Getenv("UNIX")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Shutting down server...")

		server.Shutdown(context.Background())

		if socketPath != "" {
			os.Remove(socketPath)
		}

		os.Exit(0)
	}()
	if socketPath != "" {
		defer os.Remove(socketPath)
		listener, err := net.Listen("unix", socketPath)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Listening on unix:%s", socketPath)
		log.Fatal(server.Serve(listener))
	} else {
		log.Printf("Listening on %s", server.Addr)
		log.Fatal(server.ListenAndServe())
	}
}
