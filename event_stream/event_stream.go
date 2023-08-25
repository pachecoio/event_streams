package event_stream

type Topic string

type EventStream struct {
	handlers  map[Topic][]EventHandler
	streams   map[Topic]chan Event
	listening bool
}

func NewEventStream() *EventStream {
	return &EventStream{
		handlers:  make(map[Topic][]EventHandler),
		streams:   make(map[Topic]chan Event),
		listening: false,
	}
}

func (es *EventStream) Listen() {
	if es.listening {
		return
	}

	es.listening = true
	go es.listen()
}

func (es *EventStream) listen() {
	for {
		if !es.listening {
			break
		}
		for topic, stream := range es.streams {
			select {
			case event := <-stream:
				es.handle(topic, event)
			default:
				continue
			}
		}
	}
}

func (es *EventStream) Subscribe(topic Topic, handler EventHandler) {
	if es.listening {
		panic("Event stream is already listening")
	}
	es.handlers[topic] = append(es.handlers[topic], handler)
}

func (es *EventStream) Publish(topic Topic, event Event) {
	if !es.listening {
		panic("Event stream is not listening")
	}
	if _, ok := es.streams[topic]; !ok {
		es.streams[topic] = make(chan Event)
	}

	es.streams[topic] <- event
}

func (es *EventStream) Stop() {
	es.listening = false
}

func (es *EventStream) handle(topic Topic, event Event) {
	for _, handler := range es.handlers[topic] {
		handler(event)
	}
}

type Event interface {
	Name() string
	ToJSON() ([]byte, error)
	FromJSON([]byte) error
}

type EventHandler func(event Event) error
