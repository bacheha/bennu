package organizations

import (
	"context"

	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	organizationDao *Dao
}

func NewService(organizationDao *Dao) *service {
	return &service{
		organizationDao: organizationDao,
	}
}

func (s *service) Find(ctx context.Context, filter bson.D) ([]render.Renderer, error) {
	orgs, err := s.organizationDao.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	renders := []render.Renderer{}
	for _, org := range orgs {
		renders = append(renders, org)
	}
	return renders, nil
}

func (s *service) FindById(ctx context.Context, id string) (render.Renderer, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	where := bson.D{{Key: "_id", Value: oid}}
	org, err := s.organizationDao.FindOne(ctx, where)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (s *service) Create(ctx context.Context, org *Organization) (string, error) {
	id, err := s.organizationDao.Create(ctx, org)
	if err != nil {
		return "", err
	}
	return id, nil
}
