package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CollegeBTK-Projects/my-node-go-app/db"
	"github.com/CollegeBTK-Projects/my-node-go-app/repository"
)

func main() {
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	err := db.ConnectDB(ctx)
	if err != nil {
		log.Fatal("Ошибка: ", err)
	}
	fmt.Println("Соединение прошло")

	pool, err := db.PollConns(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = repository.AboutDB(ctx, pool)
	if err != nil {
		log.Fatal(err)
	}
	userRepo := repository.NewUserRepo(pool)
	err = userRepo.CreateTable(ctx)
	if err != nil {
		log.Printf("Ошибка: ", err)
		return
	}
	coll, err := userRepo.InsertTable(ctx)
	if err != nil {
		log.Printf("Ошибка вставки: ", err)
		return
	}
	for _, v := range coll {
		fmt.Println("Вставлена строка с id: ", v)
	}

	all, err := userRepo.GetAll(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Все: ")
	for _, s := range all {
		fmt.Println(s.ID, s.Name, s.GroupName, s.Course, s.Grade)
	}

	byGroup, err := userRepo.GetFilterGroup(ctx, "478")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Сортировка по группе: ")
	for _, s := range byGroup {
		fmt.Println(s.ID, s.Name, s.GroupName, s.Course, s.Grade)
	}

	sr, err := userRepo.GetFilterGrade(ctx, "DESC")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("По DESC: ")
	for _, s := range sr {
		fmt.Println(s.ID, s.Name, s.GroupName, s.Course, s.Grade)
	}

	page, err := userRepo.GetPage(ctx, 1, 5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Страница 1: ")
	for _, s := range page {
		fmt.Println(s.ID, s.Name, s.GroupName, s.Course, s.Grade)
	}

	if err := userRepo.UpdateUser(ctx, 8, 9.5); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Студент обновлен")

	if err := userRepo.DeleteUser(ctx, 9); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Студент удален")

	repository.DemoErrors(ctx, pool)

	if err := repository.Whole(ctx, pool); err != nil {
		log.Fatal(err)
	}
	if err := repository.Stream(ctx, pool); err != nil {
		log.Fatal(err)
	}

	n, err := repository.TxDemo(ctx, pool, []int{1, 2, 1})
	fmt.Printf("с ошибкой: строк в таблице %d, ошибка: %v\n", n, err)

	n, err = repository.TxDemo(ctx, pool, []int{1, 2, 3})
	fmt.Printf("без ошибки: строк в таблице %d, ошибка: %v\n", n, err)
	defer close()
}
