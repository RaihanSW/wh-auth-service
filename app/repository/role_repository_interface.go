package repository

import (
	"context"
	"dc-transaction-service/app/model/entity"
)

type RoleRepositoryInterface interface {
	InsertRole(ctx context.Context, data entity.Role) (int, error)
	GetRoleByID(ctx context.Context, id int) (entity.Role, error)
	GetAllRoles(ctx context.Context) ([]entity.Role, error)
	UpdateRoleByID(ctx context.Context, data entity.Role) (int64, error)
	DeleteRoleByID(ctx context.Context, id int) (int64, error)
}
