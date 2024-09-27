package auth

import (
	"context"
	"errors"

	"github.com/knuls/horus/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dao struct {
	validator *validator.Validator
	tokens    *mongo.Collection
}

func (d *Dao) Find(ctx context.Context, filter bson.D) ([]*Token, error) {
	return nil, errors.New("no impl")
}

func (d *Dao) FindOne(ctx context.Context, filter bson.D) (*Token, error) {
	return nil, errors.New("no impl")
}

func (d *Dao) Create(ctx context.Context, token *Token) (string, error) {
	return "", errors.New("no impl")
}

func (d *Dao) Update(ctx context.Context, token *Token) (*Token, error) {
	return nil, errors.New("no impl")
}

func NewDao(validator *validator.Validator, tokens *mongo.Collection) *Dao {
	return &Dao{
		validator: validator,
		tokens:    tokens,
	}
}
