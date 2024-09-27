package users

import (
	"context"

	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	userDao *Dao
}

func NewService(userDao *Dao) *service {
	return &service{
		userDao: userDao,
	}
}

func (s *service) Find(ctx context.Context, filter bson.D) ([]render.Renderer, error) {
	users, err := s.userDao.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	renders := []render.Renderer{}
	for _, user := range users {
		renders = append(renders, user)
	}
	return renders, nil
}

func (s *service) FindById(ctx context.Context, id string) (render.Renderer, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	where := bson.D{{Key: "_id", Value: oid}}
	user, err := s.userDao.FindOne(ctx, where)
	if err != nil {
		return nil, err
	}
	return user, nil
}
