package repository

import (
	"context"
	"database/sql"
	"dc-transaction-service/config"
	"fmt"
	"time"
)

type transactionRepository struct {
	db             *sql.DB
	dbReplica      *sql.DB
	timeout        time.Duration
	timeoutReplica time.Duration
	serviceName    string
}

func NewTransactionRepositoryImpl(serviceName string, cfg config.DatabaseConfig, db *sql.DB, cfgReplica config.DatabaseConfig, dbReplica *sql.DB) TransactionRepositoryInterface {
	return &transactionRepository{
		db:             db,
		dbReplica:      dbReplica,
		timeout:        cfg.Timeout * time.Second,
		timeoutReplica: cfgReplica.Timeout * time.Second,
		serviceName:    serviceName,
	}
}

func (r *transactionRepository) PingDB(ctx context.Context) error {
	ctx, cancelfunc := context.WithTimeout(ctx, r.timeout)
	defer cancelfunc()

	err := r.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("repo.PingDB : errors pinging DB : %v", err)
	}
	return nil
}
