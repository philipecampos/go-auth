package models

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Username  string        `bson:"username" json:"username"`
	Email     string        `bson:"email" json:"email"`
	Password  string        `bson:"password" json:"password"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

func (u *User) Validate() error {
	errorMessages := []string{}

	if len(u.Username) < 6 {
		errorMessages = append(errorMessages, "username must be at least 6 characters long")
	}

	if len(u.Email) < 6 {
		errorMessages = append(errorMessages, "email must be at least 6 characters long")
	}

	if len(u.Password) < 6 {
		errorMessages = append(errorMessages, "email must be at least 6 characters long")
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("%s", strings.Join(errorMessages, ", "))
	}

	return nil
}

func NewUser(username, email, password string) (*User, error) {
	user := &User{
		ID:        bson.NewObjectID(),
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID.Hex(),
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
