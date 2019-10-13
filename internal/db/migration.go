package db

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log/level"
	"github.com/sgdr/wallet-service/internal/account"
	"github.com/sgdr/wallet-service/internal/currency"
	"github.com/sgdr/wallet-service/internal/logger"
)

func CreateDbStructure(ctx context.Context, db *sql.DB) error {
	log := logger.FromContext(ctx)
	level.Info(log).Log("msg", "start db structure's creation..")

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := currency.CreateTables(ctx, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := account.CreateTables(ctx, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	level.Info(log).Log("msg", "db structure's creation finished successfully")
	return nil
}
