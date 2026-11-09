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

type userRepository struct {
	db          *sql.DB
	timeout     time.Duration
	serviceName string
}

func NewUserRepositoryImpl(serviceName string, cfg config.DatabaseConfig, db *sql.DB) UserRepositoryInterface {
	return &userRepository{
		db:          db,
		timeout:     cfg.Timeout * time.Second,
		serviceName: serviceName,
	}
}

func (r *userRepository) CheckUsernameExist(ctx context.Context, username string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var uid string
	err := r.db.QueryRowContext(ctx, query.CheckUsernameExist, username).Scan(&uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("repo.CheckUsernameExist (Scan): %w", err)
	}

	return true, nil
}

func (r *userRepository) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var uid string
	err := r.db.QueryRowContext(ctx, query.CheckEmailExist, email).Scan(&uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("repo.CheckEmailExist (Scan): %w", err)
	}

	return true, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var data entity.User
	err := r.db.QueryRowContext(ctx, query.GetUserByUsername, username).Scan(
		&data.UID,
		&data.FirstName,
		&data.LastName,
		&data.Username,
		&data.Email,
		&data.RoleID,
		&data.Password,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("repo.GetUserByUsername (Scan): %w", err)
	}

	return data, nil
}

func (r *userRepository) GetUserByUID(ctx context.Context, uid string) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var data entity.User
	err := r.db.QueryRowContext(ctx, query.GetUserByUID, uid).Scan(
		&data.UID,
		&data.FirstName,
		&data.LastName,
		&data.Username,
		&data.Email,
		&data.RoleID,
		&data.Password,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("repo.GetUserByUID (Scan): %w", err)
	}

	return data, nil
}

func (r *userRepository) InsertUser(ctx context.Context, data entity.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var uid string
	err := r.db.QueryRowContext(ctx, query.InsertUser,
		data.FirstName,
		data.LastName,
		data.Username,
		data.Email,
		data.RoleID,
		data.Password,
	).Scan(&uid)
	if err != nil {
		return "", fmt.Errorf("repo.InsertUser (Scan): %w", err)
	}

	return uid, nil
}

func (r *userRepository) UpdateUserByUID(ctx context.Context, data entity.User) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	stmt, err := r.db.PrepareContext(ctx, query.UpdateUserByUID)
	if err != nil {
		return 0, fmt.Errorf("repo.UpdateUserByUID (PrepareContext): %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx,
		data.FirstName,
		data.LastName,
		data.Username,
		data.Email,
		data.RoleID,
		data.Password,
		data.UID,
	)
	if err != nil {
		return 0, fmt.Errorf("repo.UpdateUserByUID (ExecContext): %w", err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("repo.UpdateUserByUID (RowsAffected): %w", err)
	}

	return ra, nil
}

func (r *userRepository) DeleteUserByUID(ctx context.Context, uid string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.db.ExecContext(ctx, query.DeleteUserByUID, uid)
	if err != nil {
		return 0, fmt.Errorf("repo.DeleteUserByUID (ExecContext): %w", err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("repo.DeleteUserByUID (RowsAffected): %w", err)
	}

	return ra, nil
}
