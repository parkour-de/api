package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Message struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	server := NewServer()
	log.Fatal(server.ListenAndServe())
}

func NewServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	addr := ":" + port
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}

	msg := Message{Name: name, URL: r.URL.Path}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprintf(w, "%s", jsonMsg); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
