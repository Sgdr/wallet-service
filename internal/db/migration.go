package db

import (
	"context"
	"database/sql"
	"io/ioutil"

	"github.com/go-kit/kit/log/level"
	"github.com/sgdr/wallet-service/internal/logger"
)

func CreateDbStructure(ctx context.Context, db *sql.DB, sqlFilePath string) error {
	log := logger.FromContext(ctx)
	level.Info(log).Log("msg", "start db structure's creation..")
	sqlQueries, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		return err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, string(sqlQueries))
	if err := tx.Commit(); err != nil {
		return err
	}
	level.Info(log).Log("msg", "db structure's creation finished successfully")
	return nil
}
