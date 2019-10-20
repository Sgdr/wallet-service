package account

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/sgdr/wallet-service/internal/datasource"
	"log"
)

type Repository interface {
	Store(ctx context.Context, account *Account) error
	FindAll(ctx context.Context) ([]Account, error)
	GetForUpdate(ctx context.Context, owners []string, currency string) ([]Account, error)
	UpdateBalance(ctx context.Context, id uint, newBalance int64) (int64, error)
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

const getForUpdateQuery = "SELECT a.id, a.owner, a.balance, c.id, c.ticker, c.decimals " +
	"FROM accounts a " +
	"INNER JOIN currencies c ON c.id = a.currency_id " +
	"WHERE c.ticker = $1 AND a.owner = ANY($2::text[]) " +
	"FOR UPDATE;"

func (a *accountRepository) GetForUpdate(ctx context.Context, owners []string, currency string) ([]Account, error) {
	accounts := make([]Account, 0)
	tx, ok := ctx.Value(datasource.TxFieldName).(*sql.Tx)
	if !ok {
		return accounts, errors.New("not found transaction in context")
	}
	stmt, err := tx.PrepareContext(ctx, getForUpdateQuery)
	if err != nil {
		return accounts, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, currency, pq.Array(owners))
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

const updateBalancesQuery = "UPDATE accounts SET balance=$1 WHERE id = $2;"

func (a *accountRepository) UpdateBalance(ctx context.Context, id uint, newBalance int64) (int64, error) {
	tx, ok := ctx.Value(datasource.TxFieldName).(*sql.Tx)
	if !ok {
		return 0, errors.New("not found transaction in context")
	}
	result, err := tx.ExecContext(ctx, updateBalancesQuery, newBalance, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
