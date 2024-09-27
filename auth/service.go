package auth

import (
	"context"

	"github.com/knuls/bennu/app"
	"github.com/knuls/bennu/users"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/xsrftoken"
)

type service struct {
	cfg      *app.Config
	tokenDao *Dao
	userDao  *users.Dao
}

func NewService(cfg *app.Config, tokenDao *Dao, userDao *users.Dao) *service {
	return &service{
		cfg:      cfg,
		tokenDao: tokenDao,
		userDao:  userDao,
	}
}

func (s *service) GetCSRF() string {
	token := xsrftoken.Generate(s.cfg.Auth.Csrf, "", "")
	return token
}

func (s *service) Login(ctx context.Context, email string, password string) (string, error) {
	where := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "email", Value: email}},
				bson.D{{Key: "verified", Value: true}},
			},
		},
	}

	user, err := s.userDao.FindOne(ctx, where)
	if err != nil {
		return "", err
	}

	if err := user.ComparePassword(password); err != nil {
		return "", err
	}
	// TODO: create access & refresh tokens
	return "token-here", nil
}

func (s *service) Register(ctx context.Context, user *users.User) (string, error) {
	id, err := s.userDao.Create(ctx, user)
	if err != nil {
		return "", err
	}
	// TODO: create token & send verify email with token
	return id, nil
}
