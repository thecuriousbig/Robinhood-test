package domains

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `bson:"username"`
	Password     string             `bson:"password"`
	Email        string             `bson:"email"`
	ProfileImage string             `bson:"profileImage"`
	CreatedAt    time.Time          `bson:"createdAt"`
}

type RegisterRequest struct {
	Username string
	Password string
	Email    string
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	Token string
}

type CreateUserRequest struct {
	Username string
	Password string
	Email    string
}

type UpdateUserRequest struct {
	UserId       string
	ProfileImage string
}
