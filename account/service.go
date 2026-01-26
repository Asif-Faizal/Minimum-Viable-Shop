package account

import (
	"context"
	"errors"
	"time"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/util"
	"github.com/segmentio/ksuid"
)

type Service interface {
	CreateOrUpdateAccount(ctx context.Context, account *Account) (*Account, error)
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint, take uint) ([]*Account, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	Login(ctx context.Context, email string, password string) (*AuthenticatedResponse, error)
}

type AccountService struct {
	repository         Repository
	jwtSecret          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewAccountService(
	repository Repository,
	jwtSecret string,
	accessTokenExpiry time.Duration,
	refreshTokenExpiry time.Duration,
) *AccountService {
	return &AccountService{
		repository:         repository,
		jwtSecret:          jwtSecret,
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
	}
}

func (service *AccountService) CreateOrUpdateAccount(ctx context.Context, account *Account) (*Account, error) {
	id := account.ID
	if id == "" {
		id = ksuid.New().String()
	}
	hashed := ""
	if account.Password != "" {
		hash, err := util.HashPassword(account.Password)
		if err != nil {
			return nil, err
		}
		hashed = hash
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
	return newAccount, nil
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

func (service *AccountService) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := service.repository.CheckEmailExists(ctx, email)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (service *AccountService) Login(ctx context.Context, email string, password string) (*AuthenticatedResponse, error) {
	account, err := service.repository.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !util.CheckPasswordHash(password, account.Password) {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := util.GenerateToken(account.ID, account.Email, service.jwtSecret, service.accessTokenExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := util.GenerateToken(account.ID, account.Email, service.jwtSecret, service.refreshTokenExpiry)
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:           ksuid.New().String(),
		AccountID:    account.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(service.refreshTokenExpiry),
		CreatedAt:    time.Now(),
		IsRevoked:    false,
	}

	if err := service.repository.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return &AuthenticatedResponse{
		Account:      account,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
