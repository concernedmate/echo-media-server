package models

import (
	"errors"
	"media-server/configs"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Userdata struct {
	Username   string
	Password   *string
	MaxStorage int
}

func generateToken(username string, role string, secret string) (result string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                              // Subject (user identifier)
		"iss": "media-server",                        // Issuer
		"aud": role,                                  // Audience (user role)
		"exp": time.Now().Add(time.Hour * 72).Unix(), // Expiration time
		"iat": time.Now().Unix(),                     // Issued at
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
		exp, err := claims.GetExpirationTime()
		if err != nil {
			return "", err
		}
		if exp.Before(time.Now()) {
			return "", errors.New("token expired")
		}

		if username, ok := claims["sub"].(string); !ok {
			return "", errors.New("invalid jwt token username")
		} else {
			return username, nil
		}
	} else {
		return "", errors.New("invalid jwt token")
	}
}

func Auth(username string, password string, secret string) (token string, err error) {
	var data Userdata

	err = db.QueryRow(
		`SELECT username, password, max_storage FROM users WHERE username = ?`, username,
	).Scan(&data.Username, &data.Password, &data.MaxStorage)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if data.Password == nil {
		if password != configs.DEFAULT_PASSWORD {
			return "", errors.New("invalid username or password")
		}
	} else {
		if password != *data.Password {
			return "", errors.New("invalid username or password")
		}
	}

	token, err = generateToken(username, "user", secret)
	if err != nil {
		return "", err
	}

	return token, nil
}
