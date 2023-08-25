package main

import (
	"fmt"

	"github.com/pachecoio/event-streams/event_stream"
)

func main() {
	fmt.Println("Start")

	es := event_stream.NewEventStream()
	fmt.Println("Event stream created")

	fmt.Println("Subscribing")
	es.Subscribe("test", testHandler)

	es.Listen()
	fmt.Println("Listening")
	defer es.Stop()

	es.Publish("test", SampleEvent{name: "test"})

	fmt.Println("Done")
}

func testHandler(event event_stream.Event) error {
	fmt.Printf("Received event: %s\n", event.Name())
	return nil
}

type SampleEvent struct {
	name string
}

func (e SampleEvent) Name() string {
	return e.name
}

func (e SampleEvent) ToJSON() ([]byte, error) {
	return []byte(e.name), nil
}

func (e SampleEvent) FromJSON(data []byte) error {
	e.name = string(data)
	return nil
}
