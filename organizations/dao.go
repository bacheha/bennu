package organizations

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
	validator     *validator.Validator
	organizations *mongo.Collection
}

func (d *Dao) Find(ctx context.Context, filter bson.D) ([]*Organization, error) {
	var orgs []*Organization
	cursor, err := d.organizations.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return orgs, nil
		}
		return nil, err
	}
	if err = cursor.All(ctx, &orgs); err != nil {
		return nil, err
	}
	return orgs, nil
}

func (d *Dao) FindOne(ctx context.Context, filter bson.D) (*Organization, error) {
	result := d.organizations.FindOne(ctx, filter)
	err := result.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no org found")
		}
		return nil, err
	}
	var org *Organization
	if err = result.Decode(&org); err != nil {
		return nil, err
	}
	return org, nil
}

func (d *Dao) Create(ctx context.Context, org *Organization) (string, error) {
	exists, err := d.Find(ctx, bson.D{{Key: "name", Value: org.Name}})
	if err != nil {
		return "", err
	}
	if len(exists) > 0 {
		return "", errors.New("name exists")
	}
	now := time.Now()
	org.CreatedAt = now
	org.UpdatedAt = now
	if err := d.validator.ValidateStruct(org); err != nil {
		return "", err
	}
	result, err := d.organizations.InsertOne(ctx, org)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (d *Dao) Update(ctx context.Context, org *Organization) (*Organization, error) {
	return nil, errors.New("no impl")
}

func NewDao(validator *validator.Validator, organizations *mongo.Collection) *Dao {
	return &Dao{
		validator:     validator,
		organizations: organizations,
	}
}
