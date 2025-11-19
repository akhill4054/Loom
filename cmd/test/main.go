package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/akhill4054/loom/internal/broker"
	"github.com/akhill4054/loom/pkg/encode"
	"github.com/gorilla/websocket"
)

// ConnectionURL is the WebSocket server URL.
const ConnectionURL = "ws://localhost:8080/ws"

func main() {
	fmt.Println("Simulation started")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		startPublisher("Publisher 1", "foo")
		wg.Done() // Signal completion
	}()

	for i := range 5 {
		tag := fmt.Sprintf("Subscriber %d", i)
		go startSubscriber(tag, "foo")
	}

	wg.Wait()                   // Wait until goroutines are finished running
	time.Sleep(2 * time.Second) // Wait for subscribes to finish

	fmt.Println("Finished running!")
}

func startPublisher(tag string, topic string) {
	conn, _, err := websocket.DefaultDialer.Dial(ConnectionURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	for i := 1; i <= 10; i++ {
		payload := fmt.Sprintf("%s: Message Payload %d", tag, i)
		msg := encode.Message(broker.Publish, topic, []byte(payload))

		// Write once in every 2 seconds
		time.Sleep(2 * time.Second)
		err = conn.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			log.Fatal("dial:", err)
			break
		}

		fmt.Printf("\n** %s: Published ('%s') --> %s **\n\n", tag, string(payload), topic)
	}
}

func startSubscriber(tag string, topic string) {
	conn, _, err := websocket.DefaultDialer.Dial(ConnectionURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	msg := encode.Message(broker.Subscribe, topic, []byte{})
	err = conn.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		log.Fatal("dial:", err)
		return
	}

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("dial:", err)
		}
		if msgType != websocket.BinaryMessage {
			fmt.Println("Invalid message type:", msgType)
		}

		fmt.Println(fmt.Sprintf("%s: New message: %s", tag, string(msg)))
	}
}
