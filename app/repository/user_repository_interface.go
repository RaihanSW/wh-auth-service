package repository

import (
	"context"
	"dc-transaction-service/app/model/entity"
)

type UserRepositoryInterface interface {
	CheckUsernameExist(ctx context.Context, username string) (bool, error)
	CheckEmailExist(ctx context.Context, email string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserByUID(ctx context.Context, uid string) (entity.User, error)
	InsertUser(ctx context.Context, data entity.User) (string, error)
	UpdateUserByUID(ctx context.Context, data entity.User) (int64, error)
	DeleteUserByUID(ctx context.Context, uid string) (int64, error)
}
