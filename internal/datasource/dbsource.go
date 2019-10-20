package datasource

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/go-kit/kit/log/level"
	_ "github.com/lib/pq"
	"github.com/sgdr/wallet-service/internal/config"
	"github.com/sgdr/wallet-service/internal/logger"
)

const (
	connectionRetriesCount             = 5
	delayBetweenConnectionsAttemptsSec = 2
	TxFieldName                        = "transaction"
)

var (
	dataSource *sql.DB
	once       sync.Once
)

func New(ctx context.Context, cfg config.DB) (db *sql.DB, err error) {
	once.Do(func() {
		db, err = CreateDataSource(ctx, cfg)
	})
	if err != nil {
		return nil, err
	}
	dataSource = db
	return
}

func GetDataSource() *sql.DB {
	once.Do(func() {
		panic("try to get data source before it's creation!")
	})
	return dataSource
}

// Create a connection to datasource
func CreateDataSource(ctx context.Context, cfg config.DB) (db *sql.DB, err error) {
	log := logger.FromContext(ctx)
	for retry := 0; retry < connectionRetriesCount; retry++ {
		level.Info(log).Log("msg", fmt.Sprintf("connecting to postgres %s@%s... (retry %d of %d)",
			cfg.Name, cfg.Host, retry, connectionRetriesCount))
		db, err = sql.Open("postgres", fmt.Sprintf(`
			host=%s
 			port=%s 
			user=%s
			password=%s
			dbname=%s
			sslmode=disable
		`, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name))

		if err != nil {
			level.Error(log).Log("msg", fmt.Sprintf("connecting to postgres %s@%s:%s failed: %s",
				cfg.Name, cfg.Port, cfg.Host, err))
			time.Sleep(delayBetweenConnectionsAttemptsSec * time.Second)
			continue
		}
		db.SetMaxOpenConns(cfg.MaxConnections)
		if err = db.Ping(); err != nil {
			db.Close()
			time.Sleep(delayBetweenConnectionsAttemptsSec * time.Second)
			continue
		}
		level.Info(log).Log("msg", fmt.Sprintf("successfully connected to postgres %s@%s:%s", cfg.Name, cfg.Host, cfg.Port))
		return
	}
	return
}
