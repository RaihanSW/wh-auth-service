package service

import (
	"context"
	"dc-transaction-service/app/model/entity"
	"dc-transaction-service/app/model/request"
)

type AuthServiceInterface interface {
	AuthRegister(ctx context.Context, req request.RegisterRequest) (string, int, error)
	AuthLogin(ctx context.Context, req request.LoginRequest) (entity.User, entity.AccessToken, int, error)
	AuthLogout(ctx context.Context) (int, error)
	UserGetByID(ctx context.Context, uid string) (entity.User, int, error)
	UserUpdate(ctx context.Context, uid string, req request.UserUpdateRequest) (int, error)
	UserDelete(ctx context.Context, uid string) (int, error)
}
