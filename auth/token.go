package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenScope string

const (
	ACCESS         TokenScope = "ACCESS"
	REFRESH        TokenScope = "REFRESH"
	EMAIL_VERIFY   TokenScope = "EMAIL_VERIFY"
	PASSWORD_RESET TokenScope = "PASSWORD_RESET"
)

type Token struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Scope     TokenScope         `json:"scope" bson:"scope" validate:"required,oneof=ACCESS REFRESH EMAIL_VERIFY PASSWORD_RESET"`
	Token     string             `json:"token" bson:"token" validate:"required"`
	Active    bool               `json:"active" bson:"active" validate:"required"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId" validate:"required,oid"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt" validate:"required"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt" validate:"required"`
	ExpiresAt time.Time          `json:"expiresAt" bson:"expiresAt" validate:"required"`
}

func NewToken() *Token {
	return &Token{}
}
