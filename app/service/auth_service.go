package service

import (
	"context"
	"database/sql"
	"dc-transaction-service/app/model/entity"
	"dc-transaction-service/app/model/request"
	"dc-transaction-service/app/repository"
	"dc-transaction-service/config"
	"dc-transaction-service/util"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
)

const defaultRoleID = 3 // Staff

type authServiceImpl struct {
	cfg      *config.ApplicationConfig
	userRepo repository.UserRepositoryInterface
}

func NewAuthServiceImpl(
	cfg *config.ApplicationConfig,
	userRepo repository.UserRepositoryInterface) AuthServiceInterface {

	return &authServiceImpl{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

func (m *authServiceImpl) AuthRegister(ctx context.Context, req request.RegisterRequest) (string, int, error) {
	isExistUsername, err := m.userRepo.CheckUsernameExist(ctx, req.Username)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("uc.AuthRegister (CheckUsernameExist) : %v", err)
	}
	if isExistUsername {
		return "", http.StatusBadRequest, errors.New("uc.AuthRegister : this username is already associated with another account")
	}

	isExistEmail, err := m.userRepo.CheckEmailExist(ctx, req.Email)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("uc.AuthRegister (CheckEmailExist) : %v", err)
	}
	if isExistEmail {
		return "", http.StatusBadRequest, errors.New("uc.AuthRegister : this email is already associated with another account")
	}

	roleID := req.RoleID
	if roleID == 0 {
		roleID = defaultRoleID
	}

	uid, err := m.userRepo.InsertUser(ctx, entity.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		RoleID:    roleID,
		Password:  util.HashSHA256(req.Password),
	})
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("uc.AuthRegister (InsertUser) : %v", err)
	}

	return uid, http.StatusCreated, nil
}

func (m *authServiceImpl) AuthLogin(ctx context.Context, req request.LoginRequest) (entity.User, entity.AccessToken, int, error) {
	user, err := m.userRepo.GetUserByUsername(ctx, req.Username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.User{}, entity.AccessToken{}, http.StatusUnauthorized, fmt.Errorf("uc.AuthLogin : invalid username or password")
	case err != nil:
		return entity.User{}, entity.AccessToken{}, http.StatusInternalServerError, fmt.Errorf("uc.AuthLogin (GetUserByUsername) : %v", err)
	}

	if user.Password != util.HashSHA256(req.Password) {
		return entity.User{}, entity.AccessToken{}, http.StatusUnauthorized, fmt.Errorf("uc.AuthLogin : invalid username or password")
	}

	secretKey, err := base64.StdEncoding.DecodeString(m.cfg.SecretKey)
	if err != nil {
		return entity.User{}, entity.AccessToken{}, http.StatusInternalServerError, fmt.Errorf("uc.AuthLogin (DecodeString) : %v", err)
	}

	token, err := util.GenerateToken(ctx, entity.TokenRequest{
		UID:         user.UID,
		Username:    user.Username,
		TokenExp:    m.cfg.Jwt.TokenExp,
		RefTokenExp: m.cfg.Jwt.RefTokenExp,
		Key:         string(secretKey),
	})
	if err != nil {
		return entity.User{}, entity.AccessToken{}, http.StatusInternalServerError, fmt.Errorf("uc.AuthLogin (GenerateToken) : %v", err)
	}

	return user, token, http.StatusOK, nil
}

func (m *authServiceImpl) AuthLogout(ctx context.Context) (int, error) {
	// Stateless logout: the session/JWT is managed client-side.
	return http.StatusOK, nil
}

func (m *authServiceImpl) UserGetByID(ctx context.Context, uid string) (entity.User, int, error) {
	user, err := m.userRepo.GetUserByUID(ctx, uid)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.User{}, http.StatusNotFound, fmt.Errorf("uc.UserGetByID : user not found")
	case err != nil:
		return entity.User{}, http.StatusInternalServerError, fmt.Errorf("uc.UserGetByID (GetUserByUID) : %v", err)
	}

	return user, http.StatusOK, nil
}

func (m *authServiceImpl) UserUpdate(ctx context.Context, uid string, req request.UserUpdateRequest) (int, error) {
	user, err := m.userRepo.GetUserByUID(ctx, uid)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return http.StatusNotFound, fmt.Errorf("uc.UserUpdate : user not found")
	case err != nil:
		return http.StatusInternalServerError, fmt.Errorf("uc.UserUpdate (GetUserByUID) : %v", err)
	}

	if req.Username != user.Username {
		isExist, err := m.userRepo.CheckUsernameExist(ctx, req.Username)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("uc.UserUpdate (CheckUsernameExist) : %v", err)
		}
		if isExist {
			return http.StatusBadRequest, errors.New("uc.UserUpdate : this username is already associated with another account")
		}
	}

	if req.Email != user.Email {
		isExist, err := m.userRepo.CheckEmailExist(ctx, req.Email)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("uc.UserUpdate (CheckEmailExist) : %v", err)
		}
		if isExist {
			return http.StatusBadRequest, errors.New("uc.UserUpdate : this email is already associated with another account")
		}
	}

	password := user.Password
	if req.Password != "" {
		password = util.HashSHA256(req.Password)
	}

	resultRows, err := m.userRepo.UpdateUserByUID(ctx, entity.User{
		UID:       uid,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		RoleID:    req.RoleID,
		Password:  password,
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("uc.UserUpdate (UpdateUserByUID) : %v", err)
	}
	if resultRows < 1 {
		return http.StatusNotFound, fmt.Errorf("uc.UserUpdate : user not found")
	}

	return http.StatusOK, nil
}

func (m *authServiceImpl) UserDelete(ctx context.Context, uid string) (int, error) {
	resultRows, err := m.userRepo.DeleteUserByUID(ctx, uid)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("uc.UserDelete (DeleteUserByUID) : %v", err)
	}
	if resultRows < 1 {
		return http.StatusNotFound, fmt.Errorf("uc.UserDelete : user not found")
	}

	return http.StatusOK, nil
}
