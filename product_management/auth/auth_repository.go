package auth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"product_management/utils"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepositoryInterface interface {
	CreateAccount(ctx context.Context, p *Account) (*Account, map[string]interface{}, error)
	Login(ctx context.Context, email string, password string) (*Account, map[string]interface{}, time.Time, string, error)
	Refresh(ctx context.Context, cookie *http.Cookie) (map[string]interface{}, time.Time, string, error)
}

func NewAuthRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) Validate(account *Account) (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return utils.Message(400, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return utils.Message(400, "Password is required"), false
	}

	//Email must be unique
	temp := &Account{}

	// Check for errors and duplicate email IDs
	err := r.db.Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.Message(400, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return utils.Message(400, "Email address already in use by another user."), false
	}

	return utils.Message(400, "Requirement passed"), true
}

func (r *Repository) CreateAccount(ctx context.Context, account *Account) (*Account, map[string]interface{}, error) {
	if resp, ok := r.Validate(account); !ok {
		return nil, resp, errors.New("Validation error")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	r.db.Create(account)

	if account.ID <= 0 {
		return nil, utils.Message(400, "Failed to create account, connection error."), errors.New("Failed to create account, connection error.")
	}

	// Create new JWT token for the newly registered account
	tk := &Token{}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	account.Token = tokenString

	account.Password = ""

	response := utils.Message(200, "Account has been created")
	return account, response, nil
}

func (r *Repository) Login(ctx context.Context, email, password string) (*Account, map[string]interface{}, time.Time, string, error) {
	account := &Account{}
	err := r.db.Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return account, utils.Message(400, "Email address not found"), time.Now(), "", err
		}
		return account, utils.Message(400, "Connection error. Please retry"), time.Now(), "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return account, utils.Message(400, "Invalid login credentials. Please try again"), time.Now(), "", err
	}
	account.Password = ""

	expirationTime := time.Now().Add(30 * time.Minute)
	//Create JWT token
	tk := &Token{StandardClaims: jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
	}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	account.Token = tokenString
	return account, utils.Message(200, "Logged In"), expirationTime, tokenString, nil
}

func (r *Repository) Refresh(ctx context.Context, cookie *http.Cookie) (map[string]interface{}, time.Time, string, error) {
	tknStr := cookie.Value
	claims := &Token{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_PASSWORD")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return utils.Message(401, "Signature Invalid"), time.Now(), "", err
		}
		return utils.Message(400, "Invalid Request"), time.Now(), "", err
	}

	if !tkn.Valid {
		return utils.Message(401, "Unauthorized"), time.Now(), "", err
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return utils.Message(400, "Current token has more time left to expire"), time.Now(), "", err
	}

	// Create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(30 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	if err != nil {
		return utils.Message(500, "Internal Server Error"), time.Now(), "", err
	}
	return utils.Message(200, "New token generated"), expirationTime, tokenString, err
}
