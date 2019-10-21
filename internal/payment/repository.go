package payment

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/sgdr/wallet-service/internal/datasource"
)

type Repository interface {
	Store(ctx context.Context, payment *Payment) error
	FindAll(ctx context.Context, owner string) ([]Payment, error)
}

type paymentRepository struct {
	db *sql.DB
}

func NewRepository(dbSource *sql.DB) Repository {
	return &paymentRepository{db: dbSource}
}

const insertPaymentsQuery = "INSERT INTO payments (account_from_id, account_to_id, amount, currency_id, date) VALUES ($1, $2, $3, $4, now()) RETURNING id, date;"

func (a *paymentRepository) Store(ctx context.Context, payment *Payment) error {
	tx, ok := ctx.Value(datasource.TxFieldName).(*sql.Tx)
	if !ok {
		return errors.New("not found transaction in context")
	}
	stmt, err := tx.PrepareContext(ctx, insertPaymentsQuery)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	id := uint(0)
	err = stmt.QueryRowContext(ctx, payment.AccountFrom.ID, payment.AccountTo.ID, payment.Amount, payment.Currency.ID).Scan(&id, &payment.Date)
	if err != nil {
		tx.Rollback()
		return err
	}
	payment.ID = id
	return nil
}

const findAllPaymentsQuery = "SELECT p.id, a.id, a.owner, a2.id, a2.owner, p.amount, p.date, c.id, c.ticker, c.decimals " +
	"FROM payments p " +
	"INNER JOIN accounts a ON p.account_from_id = a.id " +
	"INNER JOIN accounts a2 ON a2.id = p.account_to_id " +
	"INNER JOIN currencies c ON p.currency_id = c.id " +
	"WHERE a.owner = $1 OR a2.owner = $2;"

func (a *paymentRepository) FindAll(ctx context.Context, owner string) ([]Payment, error) {
	payments := make([]Payment, 0)
	stmt, err := a.db.Prepare(findAllPaymentsQuery)
	if err != nil {
		return payments, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, owner, owner)
	if err != nil {
		return payments, err
	}
	defer rows.Close()
	for rows.Next() {
		p := new(Payment)
		err := rows.Scan(&p.ID, &p.AccountFrom.ID, &p.AccountFrom.Owner, &p.AccountTo.ID, &p.AccountTo.Owner, &p.Amount,
			&p.Date, &p.Currency.ID, &p.Currency.Ticker, &p.Currency.Decimals)
		if err != nil {
			log.Fatal(err)
		}
		payments = append(payments, *p)
	}
	if err = rows.Err(); err != nil {
		return payments, err
	}
	return payments, err
}
