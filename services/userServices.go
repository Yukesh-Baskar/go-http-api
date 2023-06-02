package user_services

import (
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/http-crud/api/database"
	"github.com/http-crud/api/helpers"
	user_model "github.com/http-crud/api/models"
	error_handler "github.com/http-crud/api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// This is a set of functions for user registration, login, retrieval, update, and deletion using
// MongoDB and Go.
var users *mongo.Collection = database.OpenCollection(*database.Client, "users")

func RegisterUser(user *user_model.User) (*mongo.InsertOneResult, *error_handler.NewError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	count, err := users.CountDocuments(ctx, bson.M{"email": user.Email})

	if err != nil {
		return nil, &error_handler.NewError{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if count > 0 {
		return nil, &error_handler.NewError{
			Error:      "user with the same name already exist",
			StatusCode: http.StatusResetContent,
		}
	}
	var validator = validator.New()
	if err = validator.Struct(user); err != nil {
		return nil, &error_handler.NewError{
			Error:      err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	hashedPass, err := helpers.HashPassword(user.Password)

	if err != nil {
		return nil, &error_handler.NewError{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	user.Password = hashedPass
	user.ConfirmPassword = hashedPass

	user.ID = primitive.NewObjectID()

	insertionResult, err := users.InsertOne(ctx, user)

	if err != nil {
		return nil, &error_handler.NewError{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return insertionResult, nil
}

func LoginUser(email, password string) (*user_model.UserLoginResponse, *error_handler.NewError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	filter := bson.M{"email": email}
	var user user_model.User

	if err := users.FindOne(ctx, filter, nil).Decode(&user); err != nil {
		return nil, &error_handler.NewError{
			Error:      err.Error(),
			StatusCode: http.StatusNotFound,
		}
	}

	var err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return nil, &error_handler.NewError{
			Error:      err.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}

	jwt, jwtErr := helpers.GenerateJWT(&user)

	jwtRes := &user_model.UserLoginResponse{
		Accesstoken: jwt,
		ID:          user.ID,
	}

	if jwtErr != nil {
		return nil, jwtErr
	}

	return jwtRes, nil
}

func GetUserById(id string) (*user_model.User, *error_handler.NewError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, &error_handler.NewError{
			Error:      err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	filter := bson.M{"_id": objId}
	var user user_model.User
	if err := users.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, &error_handler.NewError{
			Error:      err.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}

	return &user, nil
}

func UpdateUser(user *user_model.User, id string) (*mongo.UpdateResult, *error_handler.NewError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, &error_handler.NewError{
			Error:      err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	filter := bson.M{"_id": objId}

	var updateObj bson.D

	if user.Name != "" {
		updateObj = append(updateObj, bson.E{"name", user.Name})
	}

	var userData *user_model.User

	if err := users.FindOne(ctx, filter, nil).Decode(&userData); err != nil {
		return nil, &error_handler.NewError{
			Error:      err.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}

	if user.Email != "" {
		count, err := users.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			return nil, &error_handler.NewError{
				Error:      err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}

		if count > 0 && userData.ID != objId {
			return nil, &error_handler.NewError{
				Error:      "user with this email already exist.",
				StatusCode: http.StatusNotAcceptable,
			}
		}

		updateObj = append(updateObj, bson.E{"email", user.Email})
	}

	userData.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	upsert := true

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := users.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)

	if err != nil {
		return nil, &error_handler.NewError{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return result, nil
}

func DeleteUser(id string) (*mongo.DeleteResult, *error_handler.NewError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, &error_handler.NewError{
			Error:      err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	filter := bson.M{"_id": objId}

	dResult, err := users.DeleteOne(ctx, filter, nil)

	if err != nil {
		return nil, &error_handler.NewError{
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return dResult, nil
}
