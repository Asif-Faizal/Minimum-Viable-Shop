package account

import (
	"context"
)

type Service interface {
	CreateOrUpdateAccount(ctx context.Context, account *Account) (*Account, error)
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint, take uint) ([]*Account, error)
}

type AccountService struct {
	repository Repository
}

func NewAccountService(repository Repository) *AccountService {
	return &AccountService{repository: repository}
}

func (service *AccountService) CreateOrUpdateAccount(ctx context.Context, account *Account) (*Account, error) {
	if err := service.repository.CreateOrUpdateAccount(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (service *AccountService) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	account, err := service.repository.GetAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (service *AccountService) ListAccounts(ctx context.Context, skip uint, take uint) ([]*Account, error) {
	accounts, err := service.repository.ListAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
