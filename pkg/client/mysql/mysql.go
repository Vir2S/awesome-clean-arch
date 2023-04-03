package mysql

import (
	"awesome-clean-arch/internal/config"
	"awesome-clean-arch/pkg/logging"
	"awesome-clean-arch/pkg/utils"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Client interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func NewClient(ctx context.Context, maxAttempts int, sc config.StorageConfig) (*sql.DB, error) {
	logger := logging.GetLogger()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err = db.PingContext(ctx)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		logger.Fatal("error do with tries mysql")
	}

	return db, nil
}
