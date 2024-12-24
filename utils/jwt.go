package utils

import (
	"Tasks/model"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var secretWord = []byte("secret")

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	encodeToString := hex.EncodeToString(hash)
	return encodeToString, nil
}

// CheckPasswordHash проверяет, соответствует ли пароль заданному хешу
func CheckPasswordHash(password, hash string) bool {
	decodeString, err := hex.DecodeString(hash)
	if err != nil {
		// TODO
		return false
	}
	err = bcrypt.CompareHashAndPassword(decodeString, []byte(password))
	if err != nil {
		// TODO
		return false
	}
	return true
}

func GenerateAccessJWT(refreshJWT string) (string, error) {
	id, err := ValidateRefreshJWT(refreshJWT)
	if err != nil {
		return "", err
	}
	exp := time.Now().Add(15 * time.Minute)
	claims := &jwt.MapClaims{
		"exp": exp.Unix(),
		"id":  id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString(secretWord)
	if err != nil {
		return "", nil
	}
	return str, nil
}
func ValidateAccessJWT(tokenString string) (float64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing error: %v", token.Header["alg"])
		}
		return secretWord, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(int64); ok {
			if time.Unix(exp, 0).Before(time.Now()) {
				return 0, fmt.Errorf("token expired")
			}
		}
		id := claims["id"].(float64)
		return id, nil
	}

	return 0, fmt.Errorf("invalid token, login again")
}
func GenerateRefreshJWT(person *model.Person) (string, error) {
	exp := time.Now().Add(24 * time.Hour)
	claims := &jwt.MapClaims{
		"exp": exp.Unix(),
		"id":  person.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString(secretWord)
	if err != nil {
		return "", nil
	}
	return str, nil
}
func ValidateRefreshJWT(tokenString string) (float64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing error: %v", token.Header["alg"])
		}
		return secretWord, nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := claims["id"].(float64)
		return id, nil
	}
	return 0, fmt.Errorf("invalid token, login again")
}
func GetAccessFromRefresh(refresh string) (string, error) {
	_, err := ValidateRefreshJWT(refresh)
	if err != nil {
		return "", err
	}
	accessJWT, err := GenerateAccessJWT(refresh)
	if err != nil {
		return "", err
	}
	return accessJWT, nil
}
