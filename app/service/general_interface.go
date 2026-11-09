package service

import (
	"context"
)

type GeneralServiceInterface interface {
	PingDB(ctx context.Context) (int, error)
}
