package main

import (
	"fmt"

	"github.com/r3labs/sse/v2"
)

func main() {
	events := make(chan *sse.Event)
	client := sse.NewClient("http://localhost:8080/sse")

	client.SubscribeChan("payload", events)

	for {
		select {
		case event := <-events:
			fmt.Println("Event: ", event)
		}
	}
}
