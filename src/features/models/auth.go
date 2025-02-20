package models

import (
	"errors"
	"media-server/configs"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var ErrInvalidCreds = errors.New("invalid username or password")

func generateToken(username string, secret string) (result string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	})

	result, err = token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return result, nil
}

func CheckToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(configs.JWT_SECRET), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if expired, ok := claims["exp"].(time.Time); ok {
			if expired.Before(time.Now()) {
				return "", errors.New("token expired")
			}
		} else {
			return "", err
		}

		if result, ok := claims["username"].(string); ok {
			return result, nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

func Auth(username string, password string, secret string) (token string, err error) {
	// TODO
	if username != "admin" && password != configs.DEFAULT_PASSWORD {
		return "", ErrInvalidCreds
	}

	token, err = generateToken(username, secret)
	if err != nil {
		return "", err
	}

	return token, nil
}
