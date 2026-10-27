package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AboutDB(ctx context.Context, pool *pgxpool.Pool) error {
	about := [3]string{}

	err := pool.QueryRow(ctx, "SELECT split_part(version(), ',', 1);").Scan(&about[0])
	if err != nil {
		fmt.Fprint(os.Stderr, "Не удалось получить версию БД: ", err)
		return err
	}

	err = pool.QueryRow(ctx, "SELECT CURRENT_DATABASE();").Scan(&about[1])
	if err != nil {
		fmt.Fprint(os.Stderr, "Не удалось получить название БД: ", err)
		return err
	}

	err = pool.QueryRow(ctx, "SELECT CURRENT_USER").Scan(&about[2])
	if err != nil {
		fmt.Fprint(os.Stderr, "Не удалось получить пользователя БД: ", err)
		return err
	}
	for _, v := range about {
		fmt.Println(v)
	}
	return nil
}
