package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func GreatMsg(name string) {
	fmt.Printf("Привет, %s! Добро пожаловать в CLI-приложение группы 478\n", name)
}

func main() {
	name := flag.String("great", "", "вывод приветствия")
	info := flag.Bool("info", false, "информация о группе")
	flag.Parse()

	if len(os.Args) == 1 && *name == "" {
		flag.PrintDefaults()
		return
	}

	if *info {
		time := time.Now().Format("02-01-2006 15:04:05")
		fmt.Printf("Группа: 478\nСтудент: Мишкевич Максим\nЛабораторная работа: №17\nДата: %s", time)
		return
	}

	GreatMsg(*name)
}
