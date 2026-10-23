package main

import (
	"fmt"
	"time"
)

type ReaderFlow struct {
	chunks []string
	index  int
}

func NewReaderFlow() *ReaderFlow {
	return &ReaderFlow{
		chunks: []string{
			"Группа: 478",
			"Студент: Мишкевич Максим",
			"Лабораторная работа 21",
			"Тема: Потоки в Go",
		},
		index: 2,
	}
}

func main() {
	reader := NewReaderFlow()

	queue := make(chan string, reader.index)
	done := make(chan struct{})

	go func() {
		defer close(done)
		for chunk := range queue {
			fmt.Println("[Go] ", chunk)
			time.Sleep(300 * time.Millisecond)
		}
	}()

	for i := 0; i < 4; i++ {
		chunk := reader.chunks[i]

		select {
		case queue <- chunk:
			fmt.Printf("%s\n", chunk)
		default:
			fmt.Printf("Место нету: %s\n", chunk)
			queue <- chunk
			fmt.Println("Место освободилось")
		}
	}

	close(queue)
	<-done
}
