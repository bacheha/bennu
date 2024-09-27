package dao

import (
	"context"
	"errors"
	"time"

	"github.com/knuls/bennu/organizations"
	"github.com/knuls/horus/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type organizationDao struct {
	validator     *validator.Validator
	organizations *mongo.Collection
}

func NewOrganizationDao(db *mongo.Database, validator *validator.Validator) *organizationDao {
	return &organizationDao{
		validator:     validator,
		organizations: db.Collection(organizationsCollectionName),
	}
}

func (d *organizationDao) Find(ctx context.Context, filter Where) ([]*organizations.Organization, error) {
	var orgs []*organizations.Organization
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

func (d *organizationDao) FindOne(ctx context.Context, filter Where) (*organizations.Organization, error) {
	result := d.organizations.FindOne(ctx, filter)
	err := result.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no org found")
		}
		return nil, err
	}
	var org *organizations.Organization
	if err = result.Decode(&org); err != nil {
		return nil, err
	}
	return org, nil
}

func (d *organizationDao) Create(ctx context.Context, org *organizations.Organization) (string, error) {
	exists, err := d.Find(ctx, Where{{Key: "name", Value: org.Name}})
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

func (d *organizationDao) Update(ctx context.Context, org *organizations.Organization) (*organizations.Organization, error) {
	return nil, errors.New("no impl")
}
