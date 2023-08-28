package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/exp/slices"

	"github.com/javaman/go-loyality/internal/domain"
)

type postgresWithdrawRepository struct {
	db *sql.DB
}

func NewWithdrawRepository(connectionString string) *postgresWithdrawRepository {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		panic(err)
	}
	if err := createwWithdrawsTable(db); err != nil {
		panic(err)
	}
	return &postgresWithdrawRepository{db}
}

func createwWithdrawsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS withdraws (
			number text,
			login text,
			sum numeric,
			processed_at timestamp NOT NULL default now(),
			CONSTRAINT withdraws_pk PRIMARY KEY (number),
			CONSTRAINT withdraws_users_fk FOREIGN KEY (login) REFERENCES users(login)
		)
	`)
	return err
}

func (r *postgresWithdrawRepository) Insert(w *domain.Withdraw) error {
	_, err := r.db.Exec("INSERT INTO withdraws (number, login, sum) VALUES ($1, $2, $3)", w.Order, w.Login, fromJsonNumber(w.Sum))
	return err
}

func (r *postgresWithdrawRepository) SelectAll(login string) ([]*domain.Withdraw, error) {
	rows, err := r.db.Query("SELECT number, login, sum, processed_at FROM withdraws WHERE login=$1 ORDER BY processed_at", login)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	var withdraws []*domain.Withdraw

	for rows.Next() {
		o := new(domain.Withdraw)
		var x int64
		err := rows.Scan(&o.Order, &o.Login, &x, &o.ProcessedAt)
		o.Sum = toJsonNumber(x)
		if err != nil {
			return nil, err
		}
		withdraws = append(withdraws, o)
	}
	return withdraws, nil
}

func fromJsonNumber(x json.Number) int64 {
	str := strings.Split(string(x), ".")
	actualAccural, _ := strconv.Atoi(str[0])
	actualAccural *= 100
	if len(str) > 1 {
		switch len(str[1]) {
		case 1:
			x, _ := strconv.Atoi(str[1])
			actualAccural += x * 10
		default:
			x, _ := strconv.Atoi(str[1][0:2])
			actualAccural += x
		}
	}
	return int64(actualAccural)
}

func toJsonNumber(x int64) json.Number {
	a := x / 100
	b := x % 100
	var result json.Number

	if b < 10 && b > 0 {
		result = json.Number(fmt.Sprintf("%d.0%d", a, b))
	} else if slices.Contains([]int64{10, 20, 30, 40, 50, 60, 70, 80, 90}, b) {
		result = json.Number(fmt.Sprintf("%d.%d", a, b/10))
	} else if b == 0 {
		result = json.Number(fmt.Sprintf("%d", a))
	} else {
		result = json.Number(fmt.Sprintf("%d.%d", a, b))
	}

	return result
}
