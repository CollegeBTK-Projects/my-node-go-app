package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func ConnectDB(ctx context.Context) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка .env: ", err)
		return err
	}

	conn, err := pgx.Connect(ctx, os.Getenv("URL_DB"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка в подключении к БД: %v\n", err)
		os.Exit(1)
		return err
	}
	defer conn.Close(ctx)
	return nil
}

func PollConns(ctx context.Context) (*pgxpool.Pool, error) {
	connCfg, err := pgxpool.ParseConfig(os.Getenv("URL_DB"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка парсинга конфига с .env: ", err)
	}
	connCfg.MaxConns = 9
	pool, err := pgxpool.NewWithConfig(ctx, connCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания конфига БД: ", err)
		os.Exit(1)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		fmt.Fprintf(os.Stderr, "Ошибка пинга БД: ", err)
		return nil, err
	}
	fmt.Println("Всего пулов: ", connCfg.MaxConns+1)
	return pool, nil
}
