package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/knuls/bennu/users"
	"github.com/knuls/horus/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestUserDaoFind(t *testing.T) {

}

func TestUserDaoFindOne(t *testing.T) {
	v, err := validator.New()
	if err != nil {
		t.Error(err)
	}
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Parallel()

	mockUser := &users.User{
		ID:        primitive.NewObjectIDFromTimestamp(time.Now()),
		Email:     "some-email",
		FirstName: "some-first-name",
		LastName:  "some-last-name",
	}

	// tests
	cases := []struct {
		name      string
		err       bool
		responses []bson.D
	}{
		{
			name: "findOneSuccess",
			err:  false,
			responses: []bson.D{
				mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{{Key: "user", Value: mockUser}}),
			},
		},
		// {
		// 	name: "findOneErr",
		// 	responses: []bson.D{
		// 		{{Key: "ok", Value: 0}},
		// 	},
		// },
		// {
		// 	name: "findOneDecodeErr",
		// 	responses: []bson.D{
		// 		mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{{Key: "user", Value: 1}}),
		// 	},
		// },
	}

	// execute
	for _, testCase := range cases {
		mt.Run(testCase.name, func(t *mtest.T) {
			t.AddMockResponses(testCase.responses...)
			dao := NewUserDao(t.Client.Database("db"), v)
			result, err := dao.FindOne(context.Background(), Where{})
			if err != nil {
				t.Error(err)
			}
			if testCase.err {
				if err == nil {
					t.Error(err)
				}
			} else {
				if err != nil {
					t.Error(err)
				}
				fmt.Println(result)
				if result != mockUser {
					t.Fatal("user not same expected")
				}
			}
		})
	}
}
