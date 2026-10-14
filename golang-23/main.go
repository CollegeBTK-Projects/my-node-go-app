package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CollegeBTK-Projects/my-node-go-app/db"
	"github.com/CollegeBTK-Projects/my-node-go-app/models"
	"github.com/CollegeBTK-Projects/my-node-go-app/repository"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка .env")
	}
	ctx := context.Background()
	db, err := db.Connect(ctx, os.Getenv("MONGO"))
	if err != nil {
		log.Fatalf("Ошибка в подключении в бд: ", err)
	}

	fmt.Println(repository.GetVersion(ctx, db))
	dbs, _ := repository.GetListDatabases(ctx, db)
	for _, v := range dbs {
		fmt.Println(v)
	}
	colls, _ := repository.GetListCollections(ctx, db)
	for _, v := range colls {
		fmt.Println(v)
	}

	coll := db.Database("group478").Collection("students")
	userrepo := repository.NewUserRepository(coll)

	u := models.User{
		Name:      "Maksim",
		GroupName: "478",
		Course:    4,
		Grade:     8,
		CreatedAt: time.Now(),
	}
	id, err := userrepo.InsertOne(ctx, u)
	if err != nil {
		fmt.Printf("Ошибка вставки: ", err)
		return
	}
	fmt.Println("Вставлен: ", id)

	u1 := []models.User{
		{
			Name:      "Misha",
			GroupName: "478",
			Course:    4,
			Grade:     9,
			CreatedAt: time.Now(),
		},
		{
			Name:      "Nikita",
			GroupName: "478",
			Course:    4,
			Grade:     7,
			CreatedAt: time.Now(),
		},
	}

	if err := userrepo.InsertMany(ctx, u1); err != nil {
		fmt.Printf("Ошибка вставки: ", err)
		return
	}

	users, err := userrepo.FindByAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range users {
		fmt.Printf("ID: %s Имя: %s Группа: %s Курс: %d Оценка: %d\n", u.ID.Hex(), u.Name, u.GroupName, u.Course, u.Grade)
	}

	idFind, err := primitive.ObjectIDFromHex("6acfbebdc02973b932342116")
	if err != nil {
		log.Fatal(err)
	}

	user, err := userrepo.FindByID(ctx, idFind)
	if err != nil {
		log.Printf("Пользователь не найден или ошибка поиска: %v", err)
	} else {
		fmt.Printf("Найден пользователь: %+v\n", *user)
	}

	groupUsers, err := userrepo.FindByGroup(ctx, "478")
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range groupUsers {
		fmt.Printf("%s %s Оценка: %d\n",
			u.Name, u.GroupName, u.Grade)
	}

	sortedUsers, err := userrepo.SortByGrades(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range sortedUsers {
		fmt.Printf("%s / Оценка: %d\n", u.Name, u.Grade)
	}

	pageUsers, err := userrepo.FindPaginated(ctx, 1, 10)
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range pageUsers {
		fmt.Printf("%s Группа: %s Оценка: %d\n", u.Name, u.GroupName, u.Grade)
	}

	err = userrepo.UpdateGrade(ctx, idFind, 10)
	if err != nil {
		log.Println("UpdateGrade:", err)
	}

	matched, modified, err := userrepo.IncGradesByCourse(ctx, 4)
	if err != nil {
		log.Println("IncGradesByCourse:", err)
	} else {
		fmt.Println("Найдено:", matched, "Изменено:", modified)
	}

	err = userrepo.DeleteByID(ctx, idFind)
	if err != nil {
		log.Println("DeleteByID:", err)
	}

	deleted, err := userrepo.DeleteByMinGrade(ctx, 6)
	if err != nil {
		log.Println("DeleteByMinGrade:", err)
	} else {
		fmt.Println("Удалено пользователей:", deleted)
	}

	users1, err := userrepo.FindNamesAndGrades(ctx)
	if err != nil {
		log.Println("FindNamesAndGrades:", err)
	} else {
		for _, user := range users1 {
			fmt.Printf("Имя: %s, Оценка: %d\n", user.Name, user.Grade)
		}
	}

	users, err = userrepo.FindPaginated(ctx, 1, 5)
	if err != nil {
		log.Println("FindPaginated:", err)
	} else {
		fmt.Println("Первая страница:")
		for _, user := range users {
			fmt.Printf("%+v\n", user)
		}
	}

	if err := userrepo.ProcessUsers(ctx); err != nil {
		log.Println(err)
	}

	if err := userrepo.ProcessUsersInBatches(ctx); err != nil {
		log.Println(err)
	}

	averages, err := userrepo.AverageGradeByGroup(ctx)
	if err != nil {
		log.Println(err)
	} else {
		for _, item := range averages {
			fmt.Printf("Группа: %v, Средний балл: %v\n",
				item["_id"], item["average_grade"])
		}
	}

	counts, err := userrepo.CountByCourse(ctx)
	if err != nil {
		log.Println(err)
	} else {
		for _, item := range counts {
			fmt.Printf("Курс: %v, Количество: %v\n",
				item["_id"], item["count"])
		}
	}

	top, err := userrepo.Top10Students(ctx)
	if err != nil {
		log.Println(err)
	} else {
		for _, user := range top {
			fmt.Printf("%s — %d\n", user.Name, user.Grade)
		}
	}

	if err := userrepo.ReadStreaming(ctx); err != nil {
		log.Println(err)
	}

	if err := userrepo.CreateGroupIndex(ctx); err != nil {
		log.Println(err)
	}

	if err := userrepo.CreateGroupGradeIndex(ctx); err != nil {
		log.Println(err)
	}

	explain, err := userrepo.ExplainGroupQuery(ctx)
	if err != nil {
		log.Println(err)
	} else {
		data, err := json.MarshalIndent(explain, "", "  ")
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(string(data))
		}
	}

	if err := userrepo.CreateTextIndex(ctx); err != nil {
		log.Println(err)
	}

	results, err := userrepo.SearchByName(ctx, "Student")
	if err != nil {
		log.Println(err)
	} else {
		for _, item := range results {
			fmt.Printf("Имя: %v, Оценка: %v, Релевантность: %v\n",
				item["name"], item["grade"], item["score"])
		}
	}
}
