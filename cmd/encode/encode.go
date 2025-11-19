package main

import (
	"encoding/base64"
	"fmt"

	"github.com/akhill4054/loom/internal/broker"
)

func encodeMessage(action broker.Action, topic string, payload []byte) []byte {
	t := []byte(topic)

	// 1 (action) + 1 (tLen) + len(topic) + 4 (payload length) + len(payload)
	buf := make([]byte, 1+1+len(t)+4+len(payload))

	i := 0
	buf[i] = byte(action)
	i++

	buf[i] = byte(len(t))
	i++

	copy(buf[i:], t)
	i += len(t)

	// payload length (uint32 big endian)
	pLen := uint32(len(payload))
	buf[i+0] = byte(pLen >> 24)
	buf[i+1] = byte(pLen >> 16)
	buf[i+2] = byte(pLen >> 8)
	buf[i+3] = byte(pLen)
	i += 4

	copy(buf[i:], payload)

	return buf
}

func main() {
	action := broker.Publish
	topic := "foo1"
	payload := []byte("foobar")

	msg := encodeMessage(action, topic, payload)

	fmt.Println("Plain:", msg)

	fmt.Println("String:", string(msg))

	encoded := base64.StdEncoding.EncodeToString(msg)
	fmt.Println("base64:", encoded)

	decoded, _ := base64.RawStdEncoding.DecodeString(encoded)
	fmt.Println("Test base64 (decoded):", decoded)
}
