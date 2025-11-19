package broker

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// HandleWebsocket -> Handles a new WebSocket connection.
func HandleWebsocket(b *Broker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Could not upgrade", http.StatusBadRequest)
			return
		}

		event := &Event{}

		event, err = ReadMessage(conn)
		if err != nil {
			log.Println("Closing connection: read error:", err)
			conn.Close()
			return
		}

		log.Println("New event:", event)

		switch event.Action {
		case Subscribe:
			handleSubscribe(b, conn, event)
		case Publish:
			handlePublish(b, conn, event)
		default:
			log.Panicln("Unknown action:", event.Action)
			conn.Close()
		}
	}
}

func handleSubscribe(b *Broker, conn *websocket.Conn, event *Event) {
	defer conn.Close()

	topic := b.AddTopic(event.Topic)
	ch := topic.AddSubscriber()
	defer topic.RemoveSubscriber(ch)
	log.Println("Subbed to topic:", event.Topic)

	for event := range ch {
		fmt.Println("Sending <-", event)

		err := conn.WriteMessage(websocket.BinaryMessage, event.Payload)
		if err != nil {
			log.Println("Subscriber disconnected:", err)
			return
		}
	}
}

func handlePublish(b *Broker, conn *websocket.Conn, event *Event) {
	defer conn.Close()

	topic := b.AddTopic(event.Topic)

	// Publish first message
	topic.Publish(event)

	for {
		event, err := ReadMessage(conn)
		if err != nil {
			log.Println("handlePublish: Closing connection: read error", err)
			break
		}

		topic.Publish(event)
	}
}
