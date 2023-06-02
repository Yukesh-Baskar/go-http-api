package user_model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	Name            string             `json:"name"`
	Email           string             `json:"email"`
	Gender          string             `json:"gender"`
	Password        string             `json:"-"`
	ConfirmPassword string             `json:"-"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type UserJWTSigningStruct struct {
	ID primitive.ObjectID
	jwt.RegisteredClaims
}

type UserLoginResponse struct {
	Accesstoken string             `json:"accesstoken"`
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
}
