package postgres

import (
	"context"
	"database/sql"
	"os"

	"github.com/misharud/investment_portfolio/pkg/config"

	"github.com/cenkalti/backoff/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/powerman/structlog"
	"github.com/pressly/goose"
)

func Open(ctx context.Context, cfg *config.Database, logger *structlog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		return nil, errors.Wrap(err, "couldn't connect to db")
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	var pingDB backoff.Operation = func() error {
		err = db.PingContext(ctx)
		if err != nil {
			logger.Info(ctx, nil, "DB is not ready...backing off...")
			return errors.WithStack(err)
		}
		return nil
	}

	err = backoff.Retry(pingDB, backoff.WithContext(backoff.NewExponentialBackOff(), ctx))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return db, nil
}

type Repo struct {
	db     *sqlx.DB
	logger *structlog.Logger
}

func NewRepo(db *sqlx.DB, logger *structlog.Logger) *Repo {
	return &Repo{db, logger}
}

func (r *Repo) inTx(ctx context.Context, opts *sql.TxOptions, fn func(context.Context, *sqlx.Tx) error) (err error) {
	tx, err := r.db.BeginTxx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "couldn't fire tx")
	}
	defer func() {
		if p := recover(); p != nil {
			if rErr := tx.Rollback(); rErr != nil {
				r.logger.PrintErr("db error", "err", errors.WithStack(rErr))
			}
			panic(p)
		}
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				r.logger.PrintErr("db error", "err", errors.WithStack(rErr))
				err = errors.WithStack(err)
			}
		}
	}()

	err = fn(ctx, tx)

	if err != nil {
		return err
	}

	err = tx.Commit()

	return errors.Wrap(err, "couldn't commit tx")
}

func MigrationDB(db *sqlx.DB, migrationPath string) error {
	var err error

	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	current, err := goose.EnsureDBVersion(db.DB)
	if err != nil {
		return err
	}
	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	migrations, err := goose.CollectMigrations(migrationPath, current, int64(len(files)))
	if err != nil {
		return err
	}

	for i := range migrations {
		if err := migrations[i].Up(db.DB); err != nil {
			return err
		}
	}

	return nil
}
