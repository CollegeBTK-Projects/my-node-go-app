package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func GreatMsg(name string) {
	fmt.Printf("Привет, %s! Добро пожаловать в CLI-приложение группы 478\n", name)
}

func newGreetCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "greet",
		Short: "вывести приветствие",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("Ошибка: укажите имя через --name")
			}
			GreatMsg(name)
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "имя для приветствия")
	return cmd
}

func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "информация о группе",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			now := time.Now().Format("02-01-2006 15:04:05")
			fmt.Printf("Группа: 478\nСтудент: Мишкевич Максим\nЛабораторная работа: №17\nДата: %s\n", now)
		},
	}
}
