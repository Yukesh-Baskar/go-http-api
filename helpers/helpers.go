package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	user_model "github.com/http-crud/api/models"
	error_handler "github.com/http-crud/api/utils"
	"golang.org/x/crypto/bcrypt"
)

// The function takes a password string and returns a hashed version of it using bcrypt algorithm.
func HashPassword(password string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("error occured while hashing password %w", err)
	}

	return string(res), nil
}

// The function generates a JWT token for a user with a specified expiration time and secret key.
func GenerateJWT(user *user_model.User) (string, *error_handler.NewError) {
	userSigningStruct := user_model.UserJWTSigningStruct{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.ID.Hex(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userSigningStruct)
	fmt.Println(os.Getenv("JWT_SECRET_KEY"))
	jwt, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	if err != nil {
		return "", &error_handler.NewError{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return jwt, nil
}

// The function Marshal takes an interface and returns a JSON-encoded byte slice and an error.
func Marshal(a interface{}) ([]byte, error) {
	marshalledResponse, err := json.Marshal(a)

	if err != nil {
		return nil, err
	}

	return marshalledResponse, nil
}
