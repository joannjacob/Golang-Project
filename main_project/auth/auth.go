package auth

import (
	"main_project/models"
	"main_project/utils"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

var jwtKey = []byte(os.Getenv("TOKEN_PASSWORD"))

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//List of endpoints that doesn't require auth
		notAuth := []string{"/api/signup", "/api/login", "/api/token/refresh"}
		//current request path
		requestPath := r.URL.Path

		// Check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")

		// Token is missing
		if tokenHeader == "" {
			response = utils.Message(403, "Missing auth token")
			utils.Respond(w, response, 403)
			log.WithFields(log.Fields{"APIName": "JwtAuthentication"}).Error("Missing auth token")
			return
		}

		// Check if token matches format `Bearer {token-body}`
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = utils.Message(403, "Invalid auth token")
			utils.Respond(w, response, 403)
			log.WithFields(log.Fields{"APIName": "JwtAuthentication"}).Error("Invalid auth token")
			return
		}

		tokenPart := splitted[1] //
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// Malformed token
		if err != nil {
			response = utils.Message(403, "Invalid authentication token")
			utils.Respond(w, response, 403)
			log.WithFields(log.Fields{"APIName": "JwtAuthentication"}).Error("Invalid authenication token")
			return
		}
		// Token is invalid, maybe not signed on this server
		if !token.Valid {
			response = utils.Message(403, "Token is not valid.")
			utils.Respond(w, response, 403)
			log.WithFields(log.Fields{"APIName": "JwtAuthentication"}).Error("Token is not valid")
			return
		}
		r = r.WithContext(r.Context())
		next.ServeHTTP(w, r)

	})
}
