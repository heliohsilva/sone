package auth

import (
	"api/src/config"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func GenerateToken(userID int64) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissions["userID"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)

	return token.SignedString(config.Secret_key)
}

func ValidateToken(r *http.Request) error {
	tokenStr := extractToken(r)

	token, err := jwt.Parse(tokenStr, returnAuthKey)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("Invalid Token")
}

func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}

func ExtractUserID(r *http.Request) (int64, error) {
	tokenStr := extractToken(r)
	token, err := jwt.Parse(tokenStr, returnAuthKey)

	if err != nil {
		return 0, err
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userIDflt, ok := permissions["userID"].(float64); ok {
			userID := int64(userIDflt)
			return userID, nil
		}
	}

	return 0, errors.New("invalid Token")
}

func returnAuthKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Assingning method wrong %v", token.Header["alg"])
	}

	return config.Secret_key, nil
}
