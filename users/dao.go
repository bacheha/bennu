package users

import (
	"context"
	"errors"
	"time"

	"github.com/knuls/horus/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dao struct {
	validator *validator.Validator
	users     *mongo.Collection
}

func (d *Dao) Find(ctx context.Context, filter bson.D) ([]*User, error) {
	var users []*User
	cursor, err := d.users.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return users, nil
		}
		return nil, err
	}
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (d *Dao) FindOne(ctx context.Context, filter bson.D) (*User, error) {
	result := d.users.FindOne(ctx, filter)
	err := result.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no user found")
		}
		return nil, err
	}
	var user *User
	if err = result.Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (d *Dao) Create(ctx context.Context, user *User) (string, error) {
	exists, err := d.Find(ctx, bson.D{{Key: "email", Value: user.Email}})
	if err != nil {
		return "", err
	}
	if len(exists) > 0 {
		return "", errors.New("email exists")
	}
	if err := user.HashPassword(); err != nil {
		return "", err
	}
	now := time.Now()
	user.Verified = false
	user.CreatedAt = now
	user.UpdatedAt = now
	if err := d.validator.ValidateStruct(user); err != nil {
		return "", err
	}
	result, err := d.users.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (d *Dao) Update(ctx context.Context, user *User) (*User, error) {
	return nil, errors.New("no impl")
}

func NewDao(validator *validator.Validator, users *mongo.Collection) *Dao {
	return &Dao{
		validator: validator,
		users:     users,
	}
}
