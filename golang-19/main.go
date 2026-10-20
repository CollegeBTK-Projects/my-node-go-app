package main

import (
	"errors"
	"fmt"
)

type Listener func(args ...any)

type entry struct {
	id   int
	fn   Listener
	once bool
}

type EventEmitter struct {
	listeners map[string][]entry
	nextID    int
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{listeners: make(map[string][]entry)}
}

type DatabaseConnection struct {
	emitter   *EventEmitter
	connected bool
}

func NewDatabaseConnection() *DatabaseConnection {
	return &DatabaseConnection{
		emitter: NewEventEmitter(),
	}
}

func (e *EventEmitter) add(event string, l Listener, once bool) int {
	e.nextID++
	e.listeners[event] = append(e.listeners[event], entry{e.nextID, l, once})
	return e.nextID
}

func (e *EventEmitter) On(event string, l Listener) int {
	return e.add(event, l, false)
}

func (e *EventEmitter) Once(event string, l Listener) int {
	return e.add(event, l, true)
}

func (e *EventEmitter) AddListener(event string, l Listener) int {
	return e.On(event, l)
}

func (e *EventEmitter) Emit(event string, args ...any) bool {
	current := e.listeners[event]
	if len(current) == 0 {
		return false
	}

	var kept []entry
	for _, en := range current {
		if !en.once {
			kept = append(kept, en)
		}
		en.fn(args...)
	}
	e.listeners[event] = kept
	return true
}

func (e *EventEmitter) ListenerCount(event string) int {
	return len(e.listeners[event])
}

func (e *EventEmitter) Listeners(event string) []Listener {
	var res []Listener
	for _, en := range e.listeners[event] {
		res = append(res, en.fn)
	}
	return res
}

func (e *EventEmitter) RemoveListener(event string, id int) bool {
	ls := e.listeners[event]
	for i, en := range ls {
		if en.id == id {
			e.listeners[event] = append(ls[:i], ls[i+1:]...)
			return true
		}
	}
	return false
}

func (e *EventEmitter) RemoveAllListeners(event string) {
	delete(e.listeners, event)
}

func (db *DatabaseConnection) On(event string, l Listener) {
	db.emitter.On(event, l)
}

func (db *DatabaseConnection) connect() {
	db.connected = true
	db.emitter.Emit("connect", "Группа 478: подключение успешно")
}

func (db *DatabaseConnection) query(sql string) {
	if !db.connected {
		db.error("Нет подключения к базе данных")
		return
	}

	result := "Запрос выполнен успешно"

	db.emitter.Emit("query", sql, result)
}

func (db *DatabaseConnection) close() {
	db.connected = false
	db.emitter.Emit("close", "Соединение закрыто")
}

func (db *DatabaseConnection) error(message string) {
	db.emitter.Emit("error", errors.New(message))
}

func main() {
	emitter := NewEventEmitter()

	onnum, oncenum := 0, 0

	emitter.On("tick", func(args ...any) {
		onnum++
		fmt.Printf("on - tick #%v\n", args[0])
	})
	emitter.Once("tick", func(args ...any) {
		oncenum++
		fmt.Printf("once - tick #%v\n", args[0])
	})

	for i := 1; i <= 3; i++ {
		emitter.Emit("tick", i)
	}
	fmt.Printf("on: %d, once: %d\n", onnum, oncenum)

	fmt.Println("Слушателей tick:", emitter.ListenerCount("tick"))

	id1 := emitter.On("tick", func(args ...any) {
		fmt.Println("A слушатель A")
	})
	emitter.AddListener("tick", func(args ...any) {
		fmt.Println("B слушатель B")
	})

	fmt.Println("После добавления A и B:", emitter.ListenerCount("tick"))
	fmt.Println("Длина слушателей:", len(emitter.Listeners("tick")))

	fmt.Println("Удалён A:", emitter.RemoveListener("tick", id1))
	fmt.Println("После удаления A:", emitter.ListenerCount("tick"))
	emitter.Emit("tick", 4)

	emitter.RemoveAllListeners("tick")
	fmt.Println("После удаления:", emitter.ListenerCount("tick"))
	fmt.Println("Emit без слушателей:", emitter.Emit("tick", 5))

	emitter.On("order", func(args ...any) { fmt.Println("первый слушатель") })
	emitter.On("order", func(args ...any) { fmt.Println("второй слушатель") })
	emitter.On("order", func(args ...any) { fmt.Println("Группа: 478") })
	emitter.Emit("order")

	db := NewDatabaseConnection()
	db.On("connect", func(args ...any) {
		fmt.Println("CONNECT:", args[0])
	})
	db.On("query", func(args ...any) {
		fmt.Println("QUERY:", args[0])
		fmt.Println("RESULT:", args[1])
	})
	db.On("close", func(args ...any) {
		fmt.Println("CLOSE:", args[0])
	})
	db.On("error", func(args ...any) {
		fmt.Println("ERROR:", args[0])
	})
	fmt.Println("Группа 478")
	db.connect()
	db.query("SELECT * FROM users")
	db.close()
	db.query("SELECT * FROM users")
}
