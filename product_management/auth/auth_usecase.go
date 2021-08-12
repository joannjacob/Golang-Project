package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

type AuthUsecaseInterface interface {
	CreateAccount(ctx context.Context, a *Account) (*Account, map[string]interface{}, error)
	Login(ctx context.Context, email string, password string) (*Account, map[string]interface{}, time.Time, string, error)
	Refresh(ctx context.Context, cookie *http.Cookie) (map[string]interface{}, time.Time, string, error)
}

func NewAuthUsecase(db *gorm.DB, r AuthRepositoryInterface) *Usecase {
	return &Usecase{
		db:   db,
		repo: r,
	}
}

type Usecase struct {
	db   *gorm.DB
	repo AuthRepositoryInterface
}

func (u *Usecase) CreateAccount(ctx context.Context, account *Account) (*Account, map[string]interface{}, error) {

	account, msg, err := u.repo.CreateAccount(ctx, account)
	if err != nil {
		return nil, msg, err
	}
	return account, msg, nil
}

func (u *Usecase) Login(ctx context.Context, email string, password string) (*Account, map[string]interface{}, time.Time, string, error) {

	account, msg, expirationTime, token, err := u.repo.Login(ctx, email, password)
	if err != nil {
		return account, msg, time.Now(), "", err
	}
	return account, msg, expirationTime, token, err
}

func (u *Usecase) Refresh(ctx context.Context, cookie *http.Cookie) (map[string]interface{}, time.Time, string, error) {

	msg, expirationTime, token, err := u.repo.Refresh(ctx, cookie)
	if err != nil {
		return msg, time.Now(), "", err
	}
	return msg, expirationTime, token, err
}
