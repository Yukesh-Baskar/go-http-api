package user_controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	user_middleware "github.com/http-crud/api/middlewares"
	user_model "github.com/http-crud/api/models"
	user_services "github.com/http-crud/api/services"
	error_handler "github.com/http-crud/api/utils"
)

// This function handles the registration of a user and returns the result in JSON format.
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	res, err := user_services.RegisterUser(user_middleware.User)
	defer r.Body.Close()
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	user_email := r.FormValue("email")
	user_password := r.FormValue("password")
	defer r.Body.Close()

	jwt, err := user_services.LoginUser(user_email, user_password)

	defer r.Body.Close()

	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(jwt)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := user_services.GetUserById(r.URL.Query().Get("id"))
	defer r.Body.Close()

	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
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
	id := r.URL.Query().Get("id")
	res, err := user_services.UpdateUser(&user, id)

	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func DeletUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	res, err := user_services.DeleteUser(id)

	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(res)
}
