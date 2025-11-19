package main

import (
	"log"
	"net/http"

	"github.com/akhill4054/loom/internal/broker"
)

func main() {
	b := broker.NewBroker()

	http.HandleFunc("/ws", broker.HandleWebsocket(b))

	log.Println("Loom running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
