package dao

import (
	"context"
	"errors"
	"time"

	"github.com/knuls/bennu/users"
	"github.com/knuls/horus/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userDao struct {
	validator *validator.Validator
	users     *mongo.Collection
}

func NewUserDao(db *mongo.Database, validator *validator.Validator) *userDao {
	return &userDao{
		validator: validator,
		users:     db.Collection(usersCollectionName),
	}
}

func (d *userDao) Find(ctx context.Context, filter Where) ([]*users.User, error) {
	var users []*users.User
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

func (d *userDao) FindOne(ctx context.Context, filter Where) (*users.User, error) {
	result := d.users.FindOne(ctx, filter)
	err := result.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no user found")
		}
		return nil, err
	}
	var user *users.User
	if err = result.Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (d *userDao) Create(ctx context.Context, user *users.User) (string, error) {
	exists, err := d.Find(ctx, Where{{Key: "email", Value: user.Email}})
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

func (d *userDao) Update(ctx context.Context, user *users.User) (*users.User, error) {
	return nil, errors.New("no impl")
}
