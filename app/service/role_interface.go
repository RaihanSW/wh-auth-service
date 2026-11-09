package service

import (
	"context"
	"dc-transaction-service/app/model/entity"
	"dc-transaction-service/app/model/request"
)

type RoleServiceInterface interface {
	RoleCreate(ctx context.Context, req request.RoleRequest) (int, int, error)
	RoleGetByID(ctx context.Context, id int) (entity.Role, int, error)
	RoleGetAll(ctx context.Context) ([]entity.Role, int, error)
	RoleUpdate(ctx context.Context, id int, req request.RoleRequest) (int, error)
	RoleDelete(ctx context.Context, id int) (int, error)
}
