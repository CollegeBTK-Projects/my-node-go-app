package main

import (
	"fmt"
)

type Listener func(args ...any)

type Event struct {
	listeners map[string][]Listener
}

func NewEvent() *Event {
	return &Event{
		listeners: make(map[string][]Listener),
	}
}

func (e *Event) On(event string, l Listener) {
	e.listeners[event] = append(e.listeners[event], l)
}

func (e *Event) Emit(event string, args ...any) bool {
	ls := e.listeners[event]

	if len(ls) == 0 {
		return false
	}

	for _, l := range ls {
		l(args...)
	}
	return true
}

func main() {
	event := NewEvent()

	event.On("greet", func(args ...any) {
		name := "Не задан"
		if len(args) > 0 {
			name = fmt.Sprint(args[0])
		}
		fmt.Printf("greet - Привет, %s! \n", name)
	})
	event.On("info", func(args ...any) {
		fmt.Println("info - Группа: 478")
	})
	event.On("bye", func(args ...any) {
		fmt.Println("bye - Пока!")
	})

	event.Emit("greet", "Maksim")
	event.Emit("info")
	event.Emit("bye")
}
