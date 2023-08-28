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

type postgresOrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(connectionString string) *postgresOrderRepository {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		panic(err)
	}
	if err := createOrdersTable(db); err != nil {
		panic(err)
	}
	return &postgresOrderRepository{db}
}

func createOrdersTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			number text,
			login text,
			status text,
			accrural numeric,
			uploaded_at timestamp NOT NULL default now(),
			CONSTRAINT orders_pk PRIMARY KEY (number),
			CONSTRAINT users_fk FOREIGN KEY (login) REFERENCES users(login)
		)
	`)
	return err
}

func (r *postgresOrderRepository) Insert(o *domain.Order) error {

	_, err := r.db.Exec("INSERT INTO orders (number, login, status, accrural) VALUES ($1, $2, $3, $4)", o.Number, o.Login, o.Status, fromJsonNumber(o.Accrual))
	return err
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

func (r *postgresOrderRepository) Select(number string) (*domain.Order, error) {
	rows, err := r.db.Query("SELECT number, login, status, accrural, uploaded_at FROM orders WHERE number=$1 ORDER BY uploaded_at", number)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	if rows.Next() {
		o := new(domain.Order)

		var x int64

		err := rows.Scan(&o.Number, &o.Login, &o.Status, &x, &o.UploadedAt)

		o.Accrual = toJsonNumber(x)

		if err != nil {
			return nil, err
		}
		return o, nil
	} else {
		return nil, nil
	}
}

func (r *postgresOrderRepository) SelectAll(login string) ([]*domain.Order, error) {
	rows, err := r.db.Query("SELECT number, login, status, accrural, uploaded_at FROM orders WHERE login=$1", login)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	var orders []*domain.Order

	for rows.Next() {
		o := new(domain.Order)

		var x int64

		err := rows.Scan(&o.Number, &o.Login, &o.Status, &x, &o.UploadedAt)

		o.Accrual = toJsonNumber(x)

		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
