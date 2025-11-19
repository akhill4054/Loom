package main

import (
	"encoding/base64"
	"fmt"

	"github.com/akhill4054/loom/internal/broker"
	"github.com/akhill4054/loom/pkg/encode"
)

func main() {
	action := broker.Publish
	topic := "foo1"
	payload := []byte("foobar")

	msg := encode.Message(action, topic, payload)

	fmt.Println("Plain:", msg)

	fmt.Println("String:", string(msg))

	encoded := base64.StdEncoding.EncodeToString(msg)
	fmt.Println("base64:", encoded)

	decoded, _ := base64.RawStdEncoding.DecodeString(encoded)
	fmt.Println("Test base64 (decoded):", decoded)
}
