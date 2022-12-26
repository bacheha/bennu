package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/knuls/bennu/dao"
	"github.com/knuls/bennu/models"
	"github.com/knuls/horus/logger"
	"github.com/knuls/horus/res"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	logger     *logger.Logger
	daoFactory *dao.Factory
}

func (h *AuthHandler) Routes() *chi.Mux {
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

func (h *AuthHandler) CSRF(rw http.ResponseWriter, r *http.Request) {
	//
}

func (h *AuthHandler) Login(rw http.ResponseWriter, r *http.Request) {
	var payload *loginRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	defer r.Body.Close()
	if err == io.EOF {
		render.Render(rw, r, res.ErrDecode(err))
		return
	}
	if err != nil {
		render.Render(rw, r, res.ErrDecode(err))
		return
	}
	where := dao.Where{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "email", Value: payload.Email}},
				bson.D{{Key: "verified", Value: true}},
			},
		},
	}
	user, err := h.daoFactory.GetUserDao().FindOne(r.Context(), where)
	if err != nil {
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		render.Render(rw, r, res.ErrNotFound(errors.New("invalid username or password")))
		return
	}

	// TODO: create access & refresh tokens
	// TODO: set access token in resp & refresh token in cookie

	render.Status(r, http.StatusOK)
	render.Respond(rw, r, &res.JSON{"token": "token"})
}

func (h *AuthHandler) Register(rw http.ResponseWriter, r *http.Request) {
	var user *models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if err == io.EOF {
		render.Render(rw, r, res.ErrDecode(err))
		return
	}
	if err != nil {
		render.Render(rw, r, res.ErrDecode(err))
		return
	}
	now := time.Now()
	user.Verified = false
	user.CreatedAt = now
	user.UpdatedAt = now
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		render.Render(rw, r, res.Err(err, http.StatusInternalServerError))
		return
	}
	user.Password = string(bytes)
	id, err := h.daoFactory.GetUserDao().Create(r.Context(), user)
	if err != nil {
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}

	// TODO: create token & send verify email with token

	render.Status(r, http.StatusCreated)
	render.Respond(rw, r, &res.JSON{"id": id})
}

func (h *AuthHandler) ResetPassword(rw http.ResponseWriter, r *http.Request) {
	//
}

func (h *AuthHandler) VerifyEmail(rw http.ResponseWriter, r *http.Request) {
	//
}

func (h *AuthHandler) VerifyResetPassword(rw http.ResponseWriter, r *http.Request) {
	//
}

func (h *AuthHandler) TokenRefresh(rw http.ResponseWriter, r *http.Request) {
	//
}

func (h *AuthHandler) Logout(rw http.ResponseWriter, r *http.Request) {
	//
}

func NewAuthHandler(logger *logger.Logger, factory *dao.Factory) *AuthHandler {
	return &AuthHandler{
		logger:     logger,
		daoFactory: factory,
	}
}
