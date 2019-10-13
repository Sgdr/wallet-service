package account

import (
	"context"
	"database/sql"
	"log"
)

type Repository interface {
	Store(ctx context.Context, account *Account) error
	FindAll(ctx context.Context) ([]Account, error)
}

type accountRepository struct {
	db *sql.DB
}

func NewRepository(dbSource *sql.DB) Repository {
	return &accountRepository{db: dbSource}
}

const insertAccountQuery = "INSERT INTO accounts (owner, balance, currency_id) VALUES ($1, $2, $3) RETURNING id;"

func (a *accountRepository) Store(ctx context.Context, account *Account) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, insertAccountQuery)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	id := uint(0)
	err = stmt.QueryRowContext(ctx, account.Owner, account.Balance, account.Currency.ID).Scan(&id)
	if err != nil {
		tx.Rollback()
		return err
	}
	account.ID = id
	return tx.Commit()
}

const findAllAccountsQuery = "SELECT a.id, a.owner, a.balance, c.id, c.ticker, c.decimals" +
	" FROM accounts a " +
	"INNER JOIN currencies c " +
	"ON a.currency_id = c.id;"

func (a *accountRepository) FindAll(ctx context.Context) ([]Account, error) {
	accounts := make([]Account, 0)
	stmt, err := a.db.Prepare(findAllAccountsQuery)
	if err != nil {
		return accounts, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return accounts, err
	}
	defer rows.Close()
	for rows.Next() {
		a := new(Account)
		err := rows.Scan(&a.ID, &a.Owner, &a.Balance, &a.Currency.ID, &a.Currency.Ticker, &a.Currency.Decimals)
		if err != nil {
			log.Fatal(err)
		}
		accounts = append(accounts, *a)
	}
	if err = rows.Err(); err != nil {
		return accounts, err
	}
	return accounts, err
}
