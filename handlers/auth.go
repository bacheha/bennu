package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/knuls/bennu/app"
	"github.com/knuls/bennu/dao"
	"github.com/knuls/bennu/users"
	"github.com/knuls/horus/logger"
	"github.com/knuls/horus/res"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/xsrftoken"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authHandler struct {
	cfg        *app.Config
	logger     logger.Logger
	daoFactory dao.Factory
}

func NewAuthHandler(logger logger.Logger, factory dao.Factory, cfg *app.Config) *authHandler {
	return &authHandler{
		logger:     logger,
		daoFactory: factory,
		cfg:        cfg,
	}
}

func (h *authHandler) Routes() *chi.Mux {
	mux := chi.NewRouter()
	mux.Get("/csrf", h.CSRF)                     // GET /auth/csrf
	mux.Post("/login", h.Login)                  // POST /auth/login
	mux.Post("/register", h.Register)            // POST /auth/register
	mux.Post("/reset-password", h.ResetPassword) // POST /auth/reset-password
	mux.Post("/logout", h.Logout)                // POST /auth/logout
	mux.Route("/verify", func(mux chi.Router) {
		mux.Post("/email", h.VerifyEmail)                  // POST /auth/verify/email
		mux.Post("/reset-password", h.VerifyResetPassword) // POST /auth/verify/reset-password
	})
	mux.Route("/token", func(mux chi.Router) {
		mux.Post("/refresh", h.TokenRefresh) // POST /auth/token/refresh
	})
	return mux
}

func (h *authHandler) CSRF(rw http.ResponseWriter, r *http.Request) {
	token := xsrftoken.Generate(h.cfg.Auth.Csrf, "", "")
	render.Status(r, http.StatusOK)
	if err := render.Render(rw, r, &res.JSON{"token": token}); err != nil {
		h.logger.Error("failed to render", "error", err)
	}
}

func (h *authHandler) Login(rw http.ResponseWriter, r *http.Request) {
	var body *loginRequest
	err := json.NewDecoder(r.Body).Decode(body)
	defer r.Body.Close()
	if err != nil {
		if errors.Is(err, io.EOF) {
			h.logger.Error("empty request body", "error", err)
			render.Render(rw, r, res.ErrDecode(err))
			return
		}
		h.logger.Error("failed to decode request body", "error", err)
		render.Render(rw, r, res.ErrDecode(err))
		return
	}
	where := dao.Where{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "email", Value: body.Email}},
				bson.D{{Key: "verified", Value: true}},
			},
		},
	}
	user, err := h.daoFactory.GetUserDao().FindOne(r.Context(), where)
	if err != nil {
		h.logger.Error("failed to find user", "error", err)
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}
	if err := user.ComparePassword(body.Password); err != nil {
		h.logger.Error("failed to compare user password", "error", err)
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}

	// TODO: create access & refresh tokens
	// TODO: set access token in resp & refresh token in cookie

	render.Status(r, http.StatusOK)
	if err = render.Render(rw, r, &res.JSON{"token": "token"}); err != nil {
		h.logger.Error("failed to render", "error", err)
	}
}

func (h *authHandler) Register(rw http.ResponseWriter, r *http.Request) {
	user := users.NewUser()
	defer r.Body.Close()
	if err := user.FromJSON(r.Body); err != nil {
		h.logger.Error("failed to decode request body", "error", err)
		render.Render(rw, r, res.ErrDecode(err))
		return
	}
	id, err := h.daoFactory.GetUserDao().Create(r.Context(), user)
	if err != nil {
		h.logger.Error("failed to create user", "error", err)
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}

	// TODO: create token & send verify email with token

	render.Status(r, http.StatusCreated)
	if err = render.Render(rw, r, &res.JSON{"id": id}); err != nil {
		h.logger.Error("failed to render", "error", err)
	}
}

func (h *authHandler) ResetPassword(rw http.ResponseWriter, r *http.Request) {
	//
}

func (h *authHandler) VerifyEmail(rw http.ResponseWriter, r *http.Request) {
	//
}

func (h *authHandler) VerifyResetPassword(rw http.ResponseWriter, r *http.Request) {
	//
}

func (h *authHandler) TokenRefresh(rw http.ResponseWriter, r *http.Request) {
	//
}

func (h *authHandler) Logout(rw http.ResponseWriter, r *http.Request) {
	//
}
