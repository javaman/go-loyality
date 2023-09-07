package repository

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/javaman/go-loyality/internal/domain"
)

type postgresUserRepository struct {
	db *sql.DB
}

func NewUserRepository(connectionSring string) *postgresUserRepository {
	db, err := sql.Open("pgx", connectionSring)
	if err != nil {
		panic(err)
	}
	if err = createUsersTable(db); err != nil {
		panic(err)
	}
	return &postgresUserRepository{db}
}

func createUsersTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			login text,
			password text,
			CONSTRAINT users_pk PRIMARY KEY (login)
		)
	`)
	return err
}

func (r *postgresUserRepository) Insert(u *domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", u.Login, u.Password)
	var e *pgconn.PgError
	if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
		return domain.ErrorLoginExists
	}
	return err
}

func (r *postgresUserRepository) Select(login string) (*domain.User, error) {
	rows, err := r.db.Query("SELECT login, password FROM users WHERE login=$1", login)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	if rows.Next() {
		u := new(domain.User)
		err := rows.Scan(&u.Login, &u.Password)
		if err != nil {
			return nil, err
		}
		return u, nil
	} else {
		return nil, nil
	}

	
}
