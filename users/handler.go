package users

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/knuls/horus/logger"
	"github.com/knuls/horus/middlewares"
	"github.com/knuls/horus/res"
	"go.mongodb.org/mongo-driver/bson"
)

type userIDCtxKey struct{}

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
	mux.Get("/", h.Find) // GET /user
	mux.Route("/{id}", func(mux chi.Router) {
		mux.Use(middlewares.ValidateObjectID("id"))
		mux.Use(userCtx)
		mux.Get("/", h.FindById) // GET /user/:id
	})
	return mux
}

func (h *handler) Find(rw http.ResponseWriter, r *http.Request) {
	users, err := h.svc.Find(r.Context(), bson.D{})
	if err != nil {
		h.logger.Error("failed to find users", "error", err)
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	if err := render.Render(rw, r, &res.JSON{"users": users}); err != nil {
		h.logger.Error("failed to render", "error", err)
		render.Render(rw, r, res.ErrRender(err))
		return
	}
}

func (h *handler) FindById(rw http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userIDCtxKey{}).(string)
	user, err := h.svc.FindById(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to find user", "error", err)
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	if err := render.Render(rw, r, &res.JSON{"user": user}); err != nil {
		h.logger.Error("failed to render", "error", err)
		render.Render(rw, r, res.ErrRender(err))
		return
	}
}

func userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), userIDCtxKey{}, chi.URLParam(r, "id"))
		next.ServeHTTP(w, r.Clone(ctx))
	})
}
