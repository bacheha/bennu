package dao

import (
	"context"
	"errors"

	"github.com/knuls/bennu/auth"
	"github.com/knuls/horus/validator"
	"go.mongodb.org/mongo-driver/mongo"
)

type tokenDao struct {
	validator *validator.Validator
	tokens    *mongo.Collection
}

func NewTokenDao(db *mongo.Database, validator *validator.Validator) *tokenDao {
	return &tokenDao{
		validator: validator,
		tokens:    db.Collection(tokensCollectionName),
	}
}

func (d *tokenDao) Find(ctx context.Context, filter Where) ([]*auth.Token, error) {
	return nil, errors.New("no impl")
}

func (d *tokenDao) FindOne(ctx context.Context, filter Where) (*auth.Token, error) {
	return nil, errors.New("no impl")
}

func (d *tokenDao) Create(ctx context.Context, token *auth.Token) (string, error) {
	return "", errors.New("no impl")
}

func (d *tokenDao) Update(ctx context.Context, token *auth.Token) (*auth.Token, error) {
	return nil, errors.New("no impl")
}
