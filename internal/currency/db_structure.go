package currency

import (
	"context"
	"database/sql"
	"github.com/go-kit/kit/log/level"
	"github.com/sgdr/wallet-service/internal/logger"
)

const createTableQuery = `
CREATE TABLE IF NOT EXISTS currencies (
  id       serial,
  ticker   text    NOT NULL,
  decimals integer NOT NULL,
  PRIMARY KEY (id)
);`

func CreateTables(ctx context.Context, tx *sql.Tx) error {
	log := logger.FromContext(ctx)
	level.Info(log).Log("msg", "create tables for currencies")

	_, err := tx.ExecContext(ctx, createTableQuery)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
