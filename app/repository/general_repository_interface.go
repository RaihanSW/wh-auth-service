package repository

import "context"

type GeneralRepositoryInterface interface {
	PingDB(ctx context.Context) error
}
