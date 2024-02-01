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
		log.Println("\nShutting down server...")

		err := server.Shutdown(context.Background())
		if err != nil {
			log.Printf("Server stopped: %s", err.Error())
		}

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
		if err = os.Chmod(socketPath, 0666); err != nil {
			log.Printf("Could not change permissions to 0666 on unix:%s", socketPath)
		}
		log.Printf("Listening on unix:%s", socketPath)
		log.Fatal(server.Serve(listener))
	} else {
		log.Printf("Listening on %s", server.Addr)
		log.Fatal(server.ListenAndServe())
	}
}
