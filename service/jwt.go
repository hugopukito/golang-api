package service

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func generateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJWT(email string) (string, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("error in retrieving home dir")
	}

	secret_jwt, err := os.ReadFile(home + "/secrets/secret_jwt.txt")
	if err != nil {
		log.Fatalln("error in read secret file")
	}

	var mySigningKey = []byte(secret_jwt)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", errors.New("something Went Wrong: " + err.Error())
	}
	return tokenString, nil
}

func parseJwt(w http.ResponseWriter, bearerToken string) jwt.MapClaims {
	bearerToken = strings.Replace(bearerToken, "Bearer ", "", 1)

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("error in retrieving home dir")
	}

	secret_jwt, err := os.ReadFile(home + "/secrets/secret_jwt.txt")
	if err != nil {
		log.Fatalln("error in read secret file")
	}

	var mySigningKey = []byte(secret_jwt)

	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			w.WriteHeader(http.StatusBadRequest)
		}
		return mySigningKey, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims
	}
	w.WriteHeader(http.StatusBadRequest)
	return nil
}

func retrieveEmail(claims jwt.MapClaims) string {
	return claims["email"].(string)
}
