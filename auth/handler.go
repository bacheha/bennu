package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/knuls/bennu/users"
	"github.com/knuls/horus/logger"
	"github.com/knuls/horus/res"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type handler struct {
	logger *logger.Logger
	svc    *service
}

func NewHandler(logger *logger.Logger, svc *service) *handler {
	return &handler{
		logger: logger,
		svc:    svc,
	}
}

func (h *handler) Routes() *chi.Mux {
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

func (h *handler) CSRF(rw http.ResponseWriter, r *http.Request) {
	token := h.svc.GetCSRF()
	render.Status(r, http.StatusOK)
	if err := render.Render(rw, r, &res.JSON{"token": token}); err != nil {
		h.logger.Error("failed to render", "error", err)
	}
}

func (h *handler) Login(rw http.ResponseWriter, r *http.Request) {
	var body *loginRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()
	if err != nil {
		if errors.Is(err, io.EOF) {
			h.logger.Error("failed to decode empty request body", "error", err)
			render.Render(rw, r, res.ErrDecode(err))
			return
		}
		h.logger.Error("failed to decode request body", "error", err)
		render.Render(rw, r, res.ErrDecode(err))
		return
	}
	token, err := h.svc.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		h.logger.Error("failed to login user", "error", err)
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}
	// TODO: set access token in resp & refresh token in cookie
	render.Status(r, http.StatusOK)
	if err = render.Render(rw, r, &res.JSON{"token": token}); err != nil {
		h.logger.Error("failed to render", "error", err)
	}
}

func (h *handler) Register(rw http.ResponseWriter, r *http.Request) {
	user := users.NewUser()
	defer r.Body.Close()
	if err := user.FromJSON(r.Body); err != nil {
		h.logger.Error("failed to decode request body", "error", err)
		render.Render(rw, r, res.ErrDecode(err))
		return
	}
	id, err := h.svc.Register(r.Context(), user)
	if err != nil {
		h.logger.Error("failed to register user", "error", err)
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
	if err = render.Render(rw, r, &res.JSON{"id": id}); err != nil {
		h.logger.Error("failed to render", "error", err)
	}
}

func (h *handler) ResetPassword(rw http.ResponseWriter, r *http.Request) {
	//
}

func (h *handler) VerifyEmail(rw http.ResponseWriter, r *http.Request) {
	//
}

func (h *handler) VerifyResetPassword(rw http.ResponseWriter, r *http.Request) {
	//
}

func (h *handler) TokenRefresh(rw http.ResponseWriter, r *http.Request) {
	//
}

func (h *handler) Logout(rw http.ResponseWriter, r *http.Request) {
	//
}
