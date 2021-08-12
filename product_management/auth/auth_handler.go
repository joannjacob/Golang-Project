package auth

import (
	"encoding/json"
	"net/http"
	"product_management/utils"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func NewAuthHandler(db *gorm.DB, u AuthUsecaseInterface) *AuthHandler {
	return &AuthHandler{
		usecase: u,
	}
}

type AuthHandler struct {
	usecase AuthUsecaseInterface
}

func (a *AuthHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	account := &Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		utils.RespondwithJSON(w, 400, utils.Message(400, "Invalid request"))
		log.WithFields(log.Fields{"APIName": "CreateAccount", "error": err}).Error("Invalid request")
		return
	}

	data, msg, err := a.usecase.CreateAccount(r.Context(), account)
	resp := msg
	statusCode := msg["status"].(interface{}).(int)
	if err != nil {
		utils.RespondwithJSON(w, statusCode, resp)
		log.WithFields(log.Fields{"APIName": "CreateAccount", "error": err}).Error(msg["message"].(interface{}).(string))
		return
	}
	resp["data"] = data
	utils.RespondwithJSON(w, statusCode, resp)
}

func (a *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	account := &Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		utils.RespondwithJSON(w, 400, utils.Message(400, "Invalid request"))
		log.WithFields(log.Fields{"APIName": "Authenticate", "error": err}).Error("Invalid request")
		return
	}

	account, msg, expirationTime, tokenString, err := a.usecase.Login(r.Context(), account.Email, account.Password)

	resp := msg
	statusCode := msg["status"].(interface{}).(int)
	if err != nil {
		utils.RespondwithJSON(w, statusCode, resp)
		log.WithFields(log.Fields{"APIName": "Refresh", "error": err}).Error(msg["message"].(interface{}).(string))
		return
	}
	resp["data"] = account
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	utils.RespondwithJSON(w, 200, resp)
}

func (a *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			utils.RespondwithJSON(w, 401, utils.Message(401, "No cookie found"))
			log.WithFields(log.Fields{"APIName": "Refresh", "error": err}).Error("No cookie found")
			return
		}
		utils.RespondwithJSON(w, 400, utils.Message(400, "Invalid request"))
		log.WithFields(log.Fields{"APIName": "Refresh", "error": err}).Error("Invalid request")
		return
	}
	msg, expirationTime, tokenString, err := a.usecase.Refresh(r.Context(), c)
	resp := msg
	statusCode := msg["status"].(interface{}).(int)
	if err != nil {
		utils.RespondwithJSON(w, statusCode, resp)
		log.WithFields(log.Fields{"APIName": "Refresh", "error": err}).Error(msg["message"].(interface{}).(string))
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	utils.RespondwithJSON(w, statusCode, resp)
}
