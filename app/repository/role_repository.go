package repository

import (
	"context"
	"database/sql"
	"dc-transaction-service/app/model/entity"
	"dc-transaction-service/app/repository/query"
	"dc-transaction-service/config"
	"fmt"
	"time"
)

type roleRepository struct {
	db          *sql.DB
	timeout     time.Duration
	serviceName string
}

func NewRoleRepositoryImpl(serviceName string, cfg config.DatabaseConfig, db *sql.DB) RoleRepositoryInterface {
	return &roleRepository{
		db:          db,
		timeout:     cfg.Timeout * time.Second,
		serviceName: serviceName,
	}
}

func (r *roleRepository) InsertRole(ctx context.Context, data entity.Role) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var id int
	err := r.db.QueryRowContext(ctx, query.InsertRole, data.Name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repo.InsertRole (Scan): %w", err)
	}

	return id, nil
}

func (r *roleRepository) GetRoleByID(ctx context.Context, id int) (entity.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var data entity.Role
	err := r.db.QueryRowContext(ctx, query.GetRoleByID, id).Scan(
		&data.ID,
		&data.Name,
	)
	if err != nil {
		return entity.Role{}, fmt.Errorf("repo.GetRoleByID (Scan): %w", err)
	}

	return data, nil
}

func (r *roleRepository) GetAllRoles(ctx context.Context) ([]entity.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query.GetAllRoles)
	if err != nil {
		return nil, fmt.Errorf("repo.GetAllRoles (QueryContext): %w", err)
	}
	defer rows.Close()

	var list []entity.Role
	for rows.Next() {
		var data entity.Role
		if err := rows.Scan(&data.ID, &data.Name); err != nil {
			return nil, fmt.Errorf("repo.GetAllRoles (Scan): %w", err)
		}
		list = append(list, data)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.GetAllRoles (Rows): %w", err)
	}

	return list, nil
}

func (r *roleRepository) UpdateRoleByID(ctx context.Context, data entity.Role) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.ExecContext(ctx, query.UpdateRoleByID, data.Name, data.ID)
	if err != nil {
		return 0, fmt.Errorf("repo.UpdateRoleByID (ExecContext): %w", err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("repo.UpdateRoleByID (RowsAffected): %w", err)
	}

	return ra, nil
}

func (r *roleRepository) DeleteRoleByID(ctx context.Context, id int) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.ExecContext(ctx, query.DeleteRoleByID, id)
	if err != nil {
		return 0, fmt.Errorf("repo.DeleteRoleByID (ExecContext): %w", err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("repo.DeleteRoleByID (RowsAffected): %w", err)
	}

	return ra, nil
}
