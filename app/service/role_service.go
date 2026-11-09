package service

import (
	"context"
	"database/sql"
	"dc-transaction-service/app/model/entity"
	"dc-transaction-service/app/model/request"
	"dc-transaction-service/app/repository"
	"dc-transaction-service/config"
	"errors"
	"fmt"
	"net/http"
)

type roleServiceImpl struct {
	cfg      *config.ApplicationConfig
	roleRepo repository.RoleRepositoryInterface
}

func NewRoleServiceImpl(
	cfg *config.ApplicationConfig,
	roleRepo repository.RoleRepositoryInterface) RoleServiceInterface {

	return &roleServiceImpl{
		cfg:      cfg,
		roleRepo: roleRepo,
	}
}

func (m *roleServiceImpl) RoleCreate(ctx context.Context, req request.RoleRequest) (int, int, error) {
	id, err := m.roleRepo.InsertRole(ctx, entity.Role{Name: req.Name})
	if err != nil {
		return 0, http.StatusInternalServerError, fmt.Errorf("uc.RoleCreate (InsertRole) : %v", err)
	}

	return id, http.StatusCreated, nil
}

func (m *roleServiceImpl) RoleGetByID(ctx context.Context, id int) (entity.Role, int, error) {
	data, err := m.roleRepo.GetRoleByID(ctx, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.Role{}, http.StatusNotFound, fmt.Errorf("uc.RoleGetByID : role not found")
	case err != nil:
		return entity.Role{}, http.StatusInternalServerError, fmt.Errorf("uc.RoleGetByID (GetRoleByID) : %v", err)
	}

	return data, http.StatusOK, nil
}

func (m *roleServiceImpl) RoleGetAll(ctx context.Context) ([]entity.Role, int, error) {
	list, err := m.roleRepo.GetAllRoles(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("uc.RoleGetAll (GetAllRoles) : %v", err)
	}

	return list, http.StatusOK, nil
}

func (m *roleServiceImpl) RoleUpdate(ctx context.Context, id int, req request.RoleRequest) (int, error) {
	resultRows, err := m.roleRepo.UpdateRoleByID(ctx, entity.Role{ID: id, Name: req.Name})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("uc.RoleUpdate (UpdateRoleByID) : %v", err)
	}
	if resultRows < 1 {
		return http.StatusNotFound, fmt.Errorf("uc.RoleUpdate : role not found")
	}

	return http.StatusOK, nil
}

func (m *roleServiceImpl) RoleDelete(ctx context.Context, id int) (int, error) {
	resultRows, err := m.roleRepo.DeleteRoleByID(ctx, id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("uc.RoleDelete (DeleteRoleByID) : %v", err)
	}
	if resultRows < 1 {
		return http.StatusNotFound, fmt.Errorf("uc.RoleDelete : role not found")
	}

	return http.StatusOK, nil
}
