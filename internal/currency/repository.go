package currency

import (
	"context"
	"database/sql"
)

type Repository interface {
	Store(ctx context.Context, account *Currency) error
}

type currencyRepository struct {
	db *sql.DB
}

func NewRepository(dbSource *sql.DB) Repository {
	return &currencyRepository{db: dbSource}
}

const insertCurrencyQuery = "INSERT INTO currencies (ticker, decimals) VALUES ($1, $2) RETURNING id;"

func (a *currencyRepository) Store(ctx context.Context, currency *Currency) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, insertCurrencyQuery)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	if _, err := stmt.ExecContext(ctx, currency.Ticker, currency.Decimals); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
