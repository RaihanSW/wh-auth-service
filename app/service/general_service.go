package service

import (
	"context"
	"fmt"
	"net/http"
	"wh-auth-service/app/repository"
	"wh-auth-service/config"
)

type generalServiceImpl struct {
	cfg             *config.ApplicationConfig
	transactionRepo repository.TransactionRepositoryInterface
}

func NewGeneralServiceImpl(
	cfg *config.ApplicationConfig,
	transactionRepo repository.TransactionRepositoryInterface) GeneralServiceInterface {

	return &generalServiceImpl{
		cfg:             cfg,
		transactionRepo: transactionRepo,
	}
}

func (m *generalServiceImpl) PingDB(ctx context.Context) (int, error) {
	err := m.transactionRepo.PingDB(ctx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("uc.PingDB (PingDB) : %v", err)
	}
	return http.StatusOK, nil
}
