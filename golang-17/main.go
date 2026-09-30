package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

func GreatMsg(name string) {
	fmt.Printf("Привет, %s! Добро пожаловать в CLI-приложение группы 478\n", name)
}

func main() {
	myflag := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	name := myflag.String("great", "", "вывод приветствия")
	info := myflag.Bool("info", false, "информация о группе")

	if err := myflag.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		os.Exit(1)
	}

	if len(os.Args) == 1 && *name == "" {
		myflag.PrintDefaults()
		return
	}

	if *info {
		time := time.Now().Format("02-01-2006 15:04:05")
		fmt.Printf("Группа: 478\nСтудент: Мишкевич Максим\nЛабораторная работа: №17\nДата: %s", time)
		return
	}

	GreatMsg(*name)
}
