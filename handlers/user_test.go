package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/knuls/bennu/dao"
	daoMocks "github.com/knuls/bennu/dao/mocks"
	"github.com/knuls/bennu/users"
	logMocks "github.com/knuls/horus/logger/mocks"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserHandler(t *testing.T) {
	t.Parallel()

	// mocks
	logger := logMocks.NewLogger()
	factory := &daoMocks.Factory{}
	errFactory := &daoMocks.ErrFactory{}
	id := primitive.NewObjectIDFromTimestamp(time.Now())
	url := fmt.Sprintf("/%s", id.Hex())

	// tests
	cases := []struct {
		name               string
		factory            dao.Factory
		method             string
		path               string
		expectedStatusCode int
		expectedBody       []*users.User
	}{
		{
			name:               "getUser",
			factory:            factory,
			method:             http.MethodGet,
			path:               "/",
			expectedStatusCode: http.StatusOK,
			expectedBody:       daoMocks.MockUsers,
		},
		{
			name:               "getUserErr",
			factory:            errFactory,
			method:             http.MethodGet,
			path:               "/",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "getUserById",
			factory:            factory,
			method:             http.MethodGet,
			path:               url,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "getUserByIdErr",
			factory:            errFactory,
			method:             http.MethodGet,
			path:               url,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	// execute
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			// target
			handler := NewUserHandler(logger, testCase.factory)
			req := httptest.NewRequest(testCase.method, testCase.path, nil)
			rr := httptest.NewRecorder()

			// serve
			handler.Routes().ServeHTTP(rr, req)

			// assert
			res := rr.Result()
			if res.StatusCode != testCase.expectedStatusCode {
				t.Fatalf("result expected to be %d, got %d", testCase.expectedStatusCode, res.StatusCode)
			}
		})
	}
}
