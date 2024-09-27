package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/knuls/bennu/app"
	"github.com/knuls/bennu/dao"
	daoMocks "github.com/knuls/bennu/dao/mocks"
	logMocks "github.com/knuls/horus/logger/mocks"
)

func TestAuthHandler(t *testing.T) {
	t.Parallel()

	// mocks
	logger := logMocks.NewLogger()
	factory := &daoMocks.Factory{}
	// errFactory := &daoMocks.ErrFactory{}
	config := &app.Config{}
	config.Auth.Csrf = "some-csrf-key"

	// tests
	cases := []*struct {
		name               string
		factory            dao.Factory
		method             string
		path               string
		body               io.Reader
		expectedStatusCode int
	}{
		{
			name:               "getCsrf",
			factory:            factory,
			method:             http.MethodGet,
			path:               "/csrf",
			body:               nil,
			expectedStatusCode: http.StatusOK,
		},
		// {
		// 	name:               "postLogin",
		// 	factory:            factory,
		// 	method:             http.MethodPost,
		// 	path:               "/login",
		// 	body:               nil,
		// 	expectedStatusCode: http.StatusOK,
		// },
		// {
		// 	name:               "postRegister",
		// 	factory:            factory,
		// 	method:             http.MethodPost,
		// 	path:               "/register",
		// 	expectedStatusCode: http.StatusOK,
		// },
		// {
		// 	name:               "postResetPassword",
		// 	factory:            factory,
		// 	method:             http.MethodPost,
		// 	path:               "/reset-password",
		// 	expectedStatusCode: http.StatusOK,
		// },
		// {
		// 	name:               "postLogout",
		// 	factory:            factory,
		// 	method:             http.MethodPost,
		// 	path:               "/logout",
		// 	expectedStatusCode: http.StatusOK,
		// },
		// {
		// 	name:               "postVerifyEmail",
		// 	factory:            factory,
		// 	method:             http.MethodPost,
		// 	path:               "/verify/email",
		// 	expectedStatusCode: http.StatusOK,
		// },
		// {
		// 	name:               "postVerifyResetPassword",
		// 	factory:            factory,
		// 	method:             http.MethodPost,
		// 	path:               "/verify/reset-password",
		// 	expectedStatusCode: http.StatusOK,
		// },
		// {
		// 	name:               "postTokenRefresh",
		// 	factory:            factory,
		// 	method:             http.MethodPost,
		// 	path:               "/token/refresh",
		// 	expectedStatusCode: http.StatusOK,
		// },
	}

	// execute
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			// target
			handler := NewAuthHandler(logger, testCase.factory, config)
			req := httptest.NewRequest(testCase.method, testCase.path, testCase.body)
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
