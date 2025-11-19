package broker

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// ReadMessage is a util function to decode an event.
func ReadMessage(conn *websocket.Conn) (*Event, error) {
	/*
		Decodes event in the following bytes format.
			+--------+----------+--------+------------+-------------+
			| ACTION | T_LEN    | TOPIC  | P_LEN      |  PAYLOAD    |
			| (1B)   | (1–2B)   | bytes  | (4B)       |  bytes      |
			+--------+----------+--------+------------+-------------+
	*/
	msgType, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	fmt.Println("Read: msg:", msg)

	if msgType != websocket.BinaryMessage {
		return nil, fmt.Errorf("Invalid message format: %v", msgType)
	}

	if len(msg) < 2 {
		return nil, fmt.Errorf("Invalid message")
	}

	i := 0

	action := Action(msg[i])
	i++

	topicLen := int(msg[i])
	i++

	if len(msg) < i+topicLen+4 {
		return nil, fmt.Errorf("Invalid topic length")
	}

	topic := string(msg[i : i+topicLen])
	i += topicLen

	if action == Subscribe {
		// Payload is not required
		return &Event{action, topic, nil}, nil
	}

	// Fancy (and efficient) way to calculate payload length using the next 4 bytes
	payloadLen := uint32(msg[i])<<24 |
		uint32(msg[i+1])<<16 |
		uint32(msg[i+2])<<8 |
		uint32(msg[i+3])
	i += 4

	if len(msg) < i+int(payloadLen) {
		return nil, fmt.Errorf("Payload too short")
	}

	payload := msg[i : i+int(payloadLen)]

	return &Event{action, topic, payload}, nil

}
