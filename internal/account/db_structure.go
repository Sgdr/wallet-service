package account

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log/level"
	"github.com/sgdr/wallet-service/internal/logger"
)

const createTableQuery = `
CREATE TABLE IF NOT EXISTS accounts (
  id          serial,
  owner       text    NOT NULL,
  balance     bigint  NOT NULL DEFAULT 0,
  currency_id integer NOT NULL,
  PRIMARY KEY (id)
);

ALTER TABLE accounts DROP CONSTRAINT IF EXISTS accounts_to_currency_foreign_key;
ALTER TABLE accounts
  ADD CONSTRAINT accounts_to_currency_foreign_key FOREIGN KEY (currency_id)
REFERENCES currencies (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE accounts DROP CONSTRAINT IF EXISTS accounts_amount_check;
ALTER TABLE accounts
  ADD CONSTRAINT accounts_amount_check CHECK (
  balance >= 0);
`

func CreateTables(ctx context.Context, tx *sql.Tx) error {
	log := logger.FromContext(ctx)
	level.Info(log).Log("msg", "create tables for accounts")

	_, err := tx.ExecContext(ctx, createTableQuery)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
