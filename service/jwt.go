package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(email, role string) (string, error) {

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
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func ParseJwt(w http.ResponseWriter, bearerToken string) {

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
			return nil, fmt.Errorf("There was an error in parsing token.")
		}
		return mySigningKey, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		println(claims["email"].(string))

		// retourner la valeur de l'email pour filtrer sur les appels
		// du service vers la bdd
	}
}
