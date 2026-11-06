package repository

import (
	"context"
	"fmt"

	"github.com/CollegeBTK-Projects/my-node-go-app/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (u *UserRepo) CreateTable(ctx context.Context) error {
	_, err := u.pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS students (
			id SERIAL PRIMARY KEY,
			name VARCHAR(20) NOT NULL,
			group_name VARCHAR(20) NOT NULL,
			course INT NOT NULL,
			grade DECIMAL(3,2),
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`)
	if err != nil {
		return err
	}
	fmt.Println("БД students создана")
	return nil
}

func (u *UserRepo) InsertTable(ctx context.Context) ([]int, error) {
	rows, err := u.pool.Query(ctx, `INSERT INTO students (name, group_name, course, grade) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8) RETURNING id`, "Maksim", "478", 4, 8.53, "Misha", "478", 4, 9)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coll []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		coll = append(coll, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return coll, nil
}

func getUsers(rows pgx.Rows) ([]models.User, error) {
	var list []models.User
	for rows.Next() {
		var s models.User
		if err := rows.Scan(&s.ID, &s.Name, &s.GroupName, &s.Course, &s.Grade); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	defer rows.Close()
	return list, nil
}

func (u *UserRepo) GetAll(ctx context.Context) ([]models.User, error) {
	rows, err := u.pool.Query(ctx,
		`SELECT id, name, group_name, course, grade FROM students ORDER BY id`)
	if err != nil {
		return nil, err
	}
	return getUsers(rows)
}

func (u *UserRepo) GetFilterGroup(ctx context.Context, group string) ([]models.User, error) {
	rows, err := u.pool.Query(ctx,
		`SELECT id, name, group_name, course, grade FROM students WHERE group_name = $1`, group)
	if err != nil {
		return nil, err
	}
	return getUsers(rows)
}

func (u *UserRepo) GetFilterGrade(ctx context.Context, type1 string) ([]models.User, error) {
	type1 = "ASC"
	rows, err := u.pool.Query(ctx,
		`SELECT id, name, group_name, course, grade FROM students ORDER BY grade `+type1+`, id`)
	if err != nil {
		return nil, err
	}
	return getUsers(rows)
}

func (u *UserRepo) GetPage(ctx context.Context, page, pageSize int) ([]models.User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	rows, err := u.pool.Query(ctx,
		`SELECT id, name, group_name, course, grade
		 FROM students
		 ORDER BY id
		 LIMIT $1 OFFSET $2`, pageSize, offset)
	if err != nil {
		return nil, err
	}
	return getUsers(rows)
}

func (u *UserRepo) UpdateUser(ctx context.Context, id int, grade float32) error {
	com, err := u.pool.Exec(ctx, `UPDATE students SET grade = $1 WHERE id = $2`, grade, id)
	if err != nil {
		return err
	}

	if com.RowsAffected() == 0 {
		return fmt.Errorf("Студент не найден: %d", id)
	}
	return nil
}

func (u *UserRepo) DeleteUser(ctx context.Context, id int) error {
	com, err := u.pool.Exec(ctx, `DELETE FROM students WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if com.RowsAffected() == 0 {
		return fmt.Errorf("Студент не найден: %d", id)
	}
	return nil
}
