package msgqueue

type Event interface {
	EventName() string
}

type EventEmitter interface {
	Emit(event Event) error
}
