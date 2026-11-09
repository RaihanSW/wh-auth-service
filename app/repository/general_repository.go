package repository

import (
	"context"
	"database/sql"
	"dc-transaction-service/config"
	"fmt"
	"time"
)

type generalRepository struct {
	db      *sql.DB
	timeout time.Duration
}

func NewGeneralRepositoryImpl(cfg config.DatabaseConfig, db *sql.DB) GeneralRepositoryInterface {
	return &generalRepository{
		db:      db,
		timeout: cfg.Timeout * time.Second,
	}
}

func (r *generalRepository) PingDB(ctx context.Context) error {
	ctx, cancelfunc := context.WithTimeout(ctx, r.timeout)
	defer cancelfunc()

	err := r.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("repo.PingDB : errors pinging DB : %v", err)
	}
	return nil
}

func transaction(db *sql.DB, txFunc func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}
