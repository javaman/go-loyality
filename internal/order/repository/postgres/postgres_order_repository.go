package repository

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

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
	_, err := r.db.Exec("INSERT INTO orders (number, login, status, accrural) VALUES ($1, $2, $3, $4)", o.Number, o.Login, o.Status, o.Accrual)
	return err
}

func (r *postgresOrderRepository) Select(number string) (*domain.Order, error) {
	rows, err := r.db.Query("SELECT number, login, status, accrural, uploaded_at FROM orders WHERE number=$1", number)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	if rows.Next() {
		o := new(domain.Order)
		err := rows.Scan(&o.Number, &o.Login, &o.Status, &o.Accrual, &o.UploadedAt)
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
		err := rows.Scan(&o.Number, &o.Login, &o.Status, &o.Accrual, &o.UploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
