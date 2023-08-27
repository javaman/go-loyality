package repository

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

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
	_, err := r.db.Exec("INSERT INTO withdraws (number, login, sum) VALUES ($1, $2, $3)", w.Order, w.Login, w.Sum)
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
		err := rows.Scan(&o.Order, &o.Login, &o.Sum, &o.ProcessedAt)
		if err != nil {
			return nil, err
		}
		withdraws = append(withdraws, o)
	}
	return withdraws, nil
}
