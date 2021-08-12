package handlers

import (
	"encoding/json"
	"main_project/models"
	"main_project/utils"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

var jwtKey = []byte(os.Getenv("TOKEN_PASSWORD"))

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		utils.Respond(w, utils.Message(400, "Invalid request"), 400)
		log.WithFields(log.Fields{"APIName": "CreateAccount", "error": err}).Error("Invalid request")
		return
	}

	resp := account.Create()
	utils.Respond(w, resp, 200)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		utils.Respond(w, utils.Message(400, "Invalid request"), 400)
		log.WithFields(log.Fields{"APIName": "Authenticate", "error": err}).Error("Invalid request")
		return
	}

	resp, expirationTime, tokenString := models.Login(account.Email, account.Password)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	utils.Respond(w, resp, 200)
}

var Refresh = func(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			utils.Respond(w, utils.Message(401, "No cookie found"), 401)
			log.WithFields(log.Fields{"APIName": "Refresh", "error": err}).Error("No cookie found")
			return
		}
		utils.Respond(w, utils.Message(400, "Invalid request"), 400)
		log.WithFields(log.Fields{"APIName": "Refresh", "error": err}).Error("Invalid request")
		return
	}
	tknStr := c.Value
	claims := &models.Token{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			utils.Respond(w, utils.Message(401, "Signature Invalid"), 401)
			log.WithFields(log.Fields{"APIName": "Refresh", "error": err}).Error("Signature Invalid")
			return
		}
		utils.Respond(w, utils.Message(400, "Invalid Request"), 400)
		log.WithFields(log.Fields{"APIName": "Refresh", "error": err}).Error("Invalid request")
		return
	}

	if !tkn.Valid {
		utils.Respond(w, utils.Message(401, "Unauthorized"), 401)
		log.WithFields(log.Fields{"APIName": "Refresh", "error": "Unauthorized"}).Error("Unauthorized")
		return
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		utils.Respond(w, utils.Message(400, "Current token has more time left to expire"), 400)
		log.WithFields(log.Fields{"APIName": "Refresh", "error": "Current token has more time left to expire"}).Error("Current token has more time left to expire")
		return
	}

	// Create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(30 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	resp := utils.Message(200, "New token generated")
	utils.Respond(w, resp, 200)
}
