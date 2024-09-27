package dao

import (
	"github.com/knuls/bennu/auth"
	"github.com/knuls/bennu/organizations"
	"github.com/knuls/bennu/users"
	"github.com/knuls/horus/validator"
	"go.mongodb.org/mongo-driver/mongo"
)

type Factory interface {
	GetUserDao() Dao[users.User]
	GetOrganizationDao() Dao[organizations.Organization]
	GetTokenDao() Dao[auth.Token]
}

type daoFactory struct {
	userDao         *userDao
	organizationDao *organizationDao
	tokenDao        *tokenDao
}

func NewDaoFactory(db *mongo.Database, validator *validator.Validator) *daoFactory {
	return &daoFactory{
		userDao:         NewUserDao(db, validator),
		organizationDao: NewOrganizationDao(db, validator),
		tokenDao:        NewTokenDao(db, validator),
	}
}

func (f *daoFactory) GetUserDao() Dao[users.User] {
	return f.userDao
}

func (f *daoFactory) GetOrganizationDao() Dao[organizations.Organization] {
	return f.organizationDao
}

func (f *daoFactory) GetTokenDao() Dao[auth.Token] {
	return f.tokenDao
}
