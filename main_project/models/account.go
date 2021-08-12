package models

import (
	"main_project/utils"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// JWT claims struct
type Token struct {
	jwt.StandardClaims
}

//Struct to represent user account
type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

var jwtKey = []byte(os.Getenv("TOKEN_PASSWORD"))

// To validate incoming user details
func (account *Account) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return utils.Message(400, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return utils.Message(400, "Password is required"), false
	}

	//Email must be unique
	temp := &Account{}

	// Check for errors and duplicate email IDs
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.Message(400, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return utils.Message(400, "Email address already in use by another user."), false
	}

	return utils.Message(400, "Requirement passed"), true
}

func (account *Account) Create() map[string]interface{} {
	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return utils.Message(400, "Failed to create account, connection error.")
	}

	// Create new JWT token for the newly registered account
	tk := &Token{}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString(jwtKey)
	account.Token = tokenString

	account.Password = ""

	response := utils.Message(200, "Account has been created")
	response["account"] = account
	return response
}

func Login(email, password string) (map[string]interface{}, time.Time, string) {
	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Message(400, "Email address not found"), time.Now(), ""
		}
		return utils.Message(400, "Connection error. Please retry"), time.Now(), ""
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return utils.Message(400, "Invalid login credentials. Please try again"), time.Now(), ""
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

	resp := utils.Message(200, "Logged In")
	resp["account"] = account
	return resp, expirationTime, tokenString
}
