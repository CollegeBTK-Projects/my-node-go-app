package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
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
	}
}

func (r *ReaderFlow) Read(p []byte) (int, error) {
	if r.index >= len(r.chunks) {
		return 0, io.EOF
	}

	chunk := r.chunks[r.index]

	if len(p) < len(chunk) {
		return 0, fmt.Errorf("буфер мал: %d", len(chunk))
	}

	n := copy(p, chunk)
	r.index++

	return n, nil
}

func pipeCopy(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}

func pipelineCopy(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer input.Close()

	output, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer output.Close()

	if _, err := io.Copy(output, input); err != nil {
		return fmt.Errorf("ошибка копирования: %w", err)
	}

	return nil
}

func main() {
	reader := NewReaderFlow()
	buf := make([]byte, 128)

	output, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}

	for {
		n, err := reader.Read(buf)

		if n > 0 {
			chunk := string(buf[:n])
			fmt.Printf("data: %s\n", chunk)

			if _, writeErr := io.WriteString(output, chunk+"\n"); writeErr != nil {
				fmt.Println("Ошибка записи:", writeErr)
				output.Close()
				return
			}
		}

		if err == io.EOF {
			fmt.Println("Чтение завершено")
			break
		}

		if err != nil {
			fmt.Println("Ошибка чтения:", err)
			output.Close()
			return
		}
	}

	if err := output.Close(); err != nil {
		fmt.Println("Ошибка закрытия файла:", err)
		return
	}

	fmt.Println("\nДемонстрация pipe")

	if err := pipeCopy("missing.txt", "pipe-copy.txt"); err != nil {
		fmt.Println("Ошибка:", err)
	}

	fmt.Println("\nДемонстрация pipeline")

	if err := pipelineCopy("missing.txt", "pipeline-copy.txt"); err != nil {
		fmt.Println("Ошибка:", err)
	}

	fmt.Println("\nКопирование файла")

	if err := pipelineCopy("output.txt", "copy.txt"); err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Файл успешно скопирован в copy.txt")

	original, err := os.ReadFile("output.txt")
	if err != nil {
		fmt.Println("Ошибка чтения output.txt:", err)
		return
	}

	copied, err := os.ReadFile("copy.txt")
	if err != nil {
		fmt.Println("Ошибка чтения copy.txt:", err)
		return
	}

	if bytes.Equal(original, copied) {
		fmt.Println("Файлы полностью совпадают")
	} else {
		fmt.Println("Файлы отличаются")
	}

	fmt.Println("\nЦепочка преобразований")

	result := strings.ToUpper(string(copied))
	result = strings.ReplaceAll(result, "GO", "GOLANG")

	fmt.Println(result)
}
