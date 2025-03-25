package configs

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var base_url = "http://localhost:3000"

var log_stack = false
var default_password = "1234"
var jwt_secret = "SECRET"

var upload_basedir = ""

func InitConfig() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	logging, err := strconv.Atoi(os.Getenv("LOG_STACK"))
	if err != nil {
		return err
	}

	log_stack = logging != 0
	base_url = os.Getenv("BASE_URL")
	default_password = os.Getenv("DEFAULT_PASSWORD")
	upload_basedir = os.Getenv("UPLOAD_BASEDIR")

	if base_url == "" || default_password == "" || upload_basedir == "" {
		return errors.New("invalid environment variables")
	}

	priv, err := os.ReadFile(os.Getenv("PRIVATE_KEY_LOCATION"))
	if err != nil {
		return err
	}
	jwt_secret = string(priv)

	return nil
}

func LOG_STACK() bool {
	return log_stack
}

func BASE_URL() string {
	return base_url
}

func DEFAULT_PASSWORD() string {
	return default_password
}

func JWT_SECRET() string {
	return jwt_secret
}

func UPLOAD_BASEDIR() string {
	return upload_basedir
}
