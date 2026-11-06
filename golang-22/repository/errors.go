package repository

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func show(name string, err error) {
	if err == nil {
		fmt.Printf("[%s] ошибки нет\n", name)
		return
	}
	fmt.Printf("[%s] %s\n    %v\n", name, Category(err), err)
}

func tryConnect(ctx context.Context, host string, port uint16, password string) error {
	cfg, err := pgx.ParseConfig(os.Getenv("URL_DB"))
	if err != nil {
		return err
	}
	cfg.Fallbacks = nil
	cfg.Host = host
	cfg.Port = port
	if password != "" {
		cfg.Password = password
	}

	cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	conn, err := pgx.ConnectConfig(cctx, cfg)
	if err != nil {
		return err
	}
	conn.Close(ctx)
	return nil
}

func DemoErrors(ctx context.Context, pool *pgxpool.Pool) {
	show("несуществующий хост", tryConnect(ctx, "no-such-host.invalid", 5432, ""))
	show("порт без сервера", tryConnect(ctx, "127.0.0.1", 5999, ""))
	show("неверный пароль", tryConnect(ctx, "127.0.0.1", 5432, "wrong_password"))

	_, err := pool.Exec(ctx, `SELEKT 1`)
	show("синтаксис", err)

	_, err = pool.Exec(ctx, `
		CREATE TEMP TABLE t1 (id int PRIMARY KEY);
		INSERT INTO t1 VALUES (1), (1)`)
	show("дубликат", err)

	_, err = pool.Exec(ctx, `
		CREATE TEMP TABLE t2 (v int NOT NULL);
		INSERT INTO t2 VALUES (NULL)`)
	show("NULL в NOT NULL", err)

	_, err = pool.Exec(ctx, `
		CREATE TEMP TABLE a (id int PRIMARY KEY);
		CREATE TEMP TABLE b (a_id int REFERENCES a(id));
		INSERT INTO b VALUES (1)`)
	show("внешний ключ", err)
}

const bigQuery = `SELECT i, md5(i::text) FROM generate_series(1, 100000) i`

func heap() uint64 {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}

func Stream(ctx context.Context, pool *pgxpool.Pool) error {
	before := heap()

	rows, err := pool.Query(ctx, bigQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var i int
		var s string
		if err := rows.Scan(&i, &s); err != nil {
			return err
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return err
	}

	fmt.Printf("потоком: %d строк, память +%.1f МБ\n",
		count, float64(heap()-before)/1048576)
	return nil
}

func Whole(ctx context.Context, pool *pgxpool.Pool) error {
	type row struct {
		I int
		S string
	}

	before := heap()

	rows, err := pool.Query(ctx, bigQuery)
	if err != nil {
		return err
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByPos[row])
	if err != nil {
		return err
	}

	fmt.Printf("целиком: %d строк, память +%.1f МБ\n",
		len(list), float64(heap()-before)/1048576)
	runtime.KeepAlive(list)
	return nil
}

func TxDemo(ctx context.Context, pool *pgxpool.Pool, ids []int) (int, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, `DROP TABLE IF EXISTS tx_t`); err != nil {
		return 0, err
	}
	if _, err := conn.Exec(ctx, `CREATE TEMP TABLE tx_t (id int PRIMARY KEY)`); err != nil {
		return 0, err
	}

	err = func() error {
		tx, err := conn.Begin(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback(ctx)

		for _, id := range ids {
			if _, err := tx.Exec(ctx, `INSERT INTO tx_t VALUES ($1)`, id); err != nil {
				return err
			}
		}
		return tx.Commit(ctx)
	}()

	var n int
	if cerr := conn.QueryRow(ctx, `SELECT count(*) FROM tx_t`).Scan(&n); cerr != nil {
		return 0, cerr
	}
	return n, err
}
