package repository

import (
	"context"
)

type TransactionRepositoryInterface interface {
	PingDB(ctx context.Context) error
}
