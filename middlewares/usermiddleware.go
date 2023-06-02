package user_middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v4"
	"github.com/http-crud/api/helpers"
	user_model "github.com/http-crud/api/models"
	error_handler "github.com/http-crud/api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var User *user_model.User

func RegisterUserMiddleware(next http.Handler) http.Handler {
	capturedErrors := []error_handler.NewError{}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" {
			eMessage, _ := json.Marshal(fmt.Sprintf("invalid method: %v", r.Method))
			w.Write(eMessage)
			return
		}

		defer r.Body.Close()

		if err := r.ParseForm(); err != nil {
			errJson, _ := json.Marshal(err)
			w.Write(errJson)
		}

		User = &user_model.User{
			Name:            r.FormValue("name"),
			Email:           r.FormValue("email"),
			Gender:          r.FormValue("gender"),
			Password:        r.FormValue("password"),
			ConfirmPassword: r.FormValue("confirmpassword"),
		}

		if strings.TrimSpace(string(User.Name)) == "" {
			capturedErrors = append(capturedErrors, error_handler.NewError{Error: "name can't be empty", StatusCode: http.StatusPartialContent})
		}
		if strings.TrimSpace(string(User.Email)) == "" {
			capturedErrors = append(capturedErrors, error_handler.NewError{Error: "email can't be empty", StatusCode: http.StatusPartialContent})
		}

		_, err := mail.ParseAddress(User.Email)

		if err != nil {
			capturedErrors = append(capturedErrors, error_handler.NewError{Error: err.Error(), StatusCode: http.StatusBadRequest})
		}

		if strings.TrimSpace(string(User.Gender)) == "" || strings.TrimSpace(string(User.Gender)) != "Male" && strings.TrimSpace(string(User.Gender)) != "Female" && strings.TrimSpace(string(User.Gender)) != "Transgender" {
			capturedErrors = append(capturedErrors, error_handler.NewError{Error: "gender can't be empty or invalid gender", StatusCode: http.StatusPartialContent})
		}

		hasNum, hasHupper, hasSpecial := verifyPassword(string(User.Password))

		if !hasNum || !hasHupper || !hasSpecial {
			capturedErrors = append(capturedErrors, error_handler.NewError{Error: fmt.Sprintf("Password missing field. hasNum: %v, hasUpper: %v, hasSpecial: %v", hasNum, hasHupper, hasSpecial)})
		}

		if strings.TrimSpace(string(User.Password)) != strings.TrimSpace(User.ConfirmPassword) {
			capturedErrors = append(capturedErrors, error_handler.NewError{Error: "password & conform password is not matched", StatusCode: http.StatusPartialContent})
		}
		if strings.TrimSpace(string(User.Password)) == "" {
			capturedErrors = append(capturedErrors, error_handler.NewError{Error: "password can't be empty", StatusCode: http.StatusPartialContent})
		}

		if len(capturedErrors) != 0 {
			byteErr, _ := json.Marshal(capturedErrors)
			w.Write([]byte(byteErr))
			capturedErrors = nil
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LoginUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" {
			e := error_handler.NewError{
				Error:      fmt.Sprintf("invalid method: %v", r.Method),
				StatusCode: http.StatusBadRequest,
			}
			eMessage, _ := json.Marshal(e)
			w.Write(eMessage)
			return
		}

		defer r.Body.Close()

		if err := r.ParseForm(); err != nil {
			errJson, _ := json.Marshal(err)
			w.Write(errJson)
		}

		user_email := r.FormValue("email")
		user_password := r.FormValue("password")

		if strings.TrimSpace(string(user_email)) == "" || strings.TrimSpace(string(user_password)) == "" {
			errMes := error_handler.NewError{
				Error:      "password or email can't be empty",
				StatusCode: http.StatusBadRequest,
			}
			jErr, _ := json.Marshal(errMes)
			w.Write(jErr)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")

		// if r.Method != "GET" || {
		// 	json.NewEncoder(w).Encode(fmt.Sprintf("invalid method: %v", r.Method))
		// 	return
		// }

		tokenString := r.Header.Get("Authorization")

		if strings.TrimSpace(tokenString) == "" {
			e, _ := helpers.Marshal(error_handler.NewError{
				Error:      "token not found",
				StatusCode: http.StatusNotFound,
			})
			w.Write(e)
			return
		}
		// This code is parsing a JWT token string and verifying its signature using a secret key.
		token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil {
			e, _ := helpers.Marshal(error_handler.NewError{
				Error:      err.Error(),
				StatusCode: http.StatusInternalServerError,
			})
			w.Write(e)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				e, _ := helpers.Marshal(error_handler.NewError{
					Error:      "jwt expired",
					StatusCode: http.StatusUnauthorized,
				})
				w.Write(e)
				return
			}

			id := r.URL.Query().Get("id")

			if !primitive.IsValidObjectID(id) {
				e, _ := json.Marshal(error_handler.NewError{
					Error:      "invalid object id",
					StatusCode: http.StatusBadRequest,
				})
				w.Write(e)
				return
			}
			if id != claims["ID"] {
				e, _ := json.Marshal(error_handler.NewError{
					Error:      "jwt not valid",
					StatusCode: http.StatusBadRequest,
				})
				w.Write(e)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}

func UpdateUserMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user user_model.User
		defer r.Body.Close()

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			fmt.Println("avgws")
			mErr, _ := json.Marshal(error_handler.NewError{
				Error:      "no data found",
				StatusCode: http.StatusNoContent,
			})
			w.Write(mErr)
			return
		}

		if user.Email != "" {
			_, err := mail.ParseAddress(user.Email)
			if err != nil {
				json.NewEncoder(w).Encode(err)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
}

func verifyPassword(s string) (hasNum, hasHupper, hasSpecial bool) {
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			hasNum = true
		case unicode.IsUpper(c):
			hasHupper = true
		case unicode.IsSymbol(c) || unicode.IsPunct(c):
			hasSpecial = true
		}
	}
	return hasNum, hasHupper, hasSpecial
}
