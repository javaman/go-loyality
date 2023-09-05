package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/exp/slices"

	"github.com/javaman/go-loyality/internal/domain"
)

type postgresBalanceRepository struct {
	db *sql.DB
}

func NewBalanceRepository(connectionSring string) *postgresBalanceRepository {
	db, err := sql.Open("pgx", connectionSring)
	if err != nil {
		panic(err)
	}
	return &postgresBalanceRepository{db}
}

func (r *postgresBalanceRepository) Select(login string) (domain.Balance, error) {
	rows, err := r.db.Query(`
		SELECT (A.DEBET - B.CREDIT), B.CREDIT
		  FROM 
		  	(SELECT coalesce(SUM(O.accrual), 0) DEBET FROM ORDERS O WHERE O.LOGIN = $1) A, 
			(SELECT coalesce(SUM(W.sum), 0) CREDIT FROM withdraws W WHERE W.LOGIN = $1) B
	`, login)

	if err != nil {
		return domain.Balance{}, err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	if rows.Next() {
		u := domain.Balance{}
		var current int64
		var withdrawn int64
		err := rows.Scan(&current, &withdrawn)
		if err != nil {
			return domain.Balance{}, err
		}
		u.Current = toJSONNumber(current)
		u.Withdrawn = toJSONNumber(withdrawn)
		return u, nil
	} 

	return domain.Balance{}, nil
}

func toJSONNumber(x int64) json.Number {
	isNegative := x < 0
	if isNegative {
		x = -x
	}

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

	if isNegative {
		result = "-" + result
	}

	return result
}
