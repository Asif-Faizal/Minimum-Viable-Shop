package account

import (
	"context"

	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
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
	id := account.ID
	if id == "" {
		id = ksuid.New().String()
	}
	hashed := ""
	if account.Password != "" {
		b, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashed = string(b)
	}
	newAccount := &Account{
		ID:       id,
		Name:     account.Name,
		Email:    account.Email,
		Password: hashed,
	}
	if _, err := service.repository.CreateOrUpdateAccount(ctx, newAccount); err != nil {
		return nil, err
	}
	return account, nil
}

func (service *AccountService) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	account, err := service.repository.GetAccountById(ctx, id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (service *AccountService) ListAccounts(ctx context.Context, skip uint, take uint) ([]*Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	accounts, err := service.repository.ListAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
