package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/knuls/bennu/dao"
	"github.com/knuls/horus/logger"
	"github.com/knuls/horus/middlewares"
	"github.com/knuls/horus/res"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userIDCtxKey struct{}

type UserHandler struct {
	logger     *logger.Logger
	daoFactory *dao.Factory
}

func (h *UserHandler) Routes() *chi.Mux {
	mux := chi.NewRouter()
	mux.Get("/", h.Find) // GET /user
	mux.Route("/{id}", func(mux chi.Router) {
		mux.Use(middlewares.ValidateObjectID("id"))
		mux.Use(UserCtx)
		mux.Get("/", h.FindById) // GET /user/:id
	})
	return mux
}

func (h *UserHandler) Find(rw http.ResponseWriter, r *http.Request) {
	users, err := h.daoFactory.GetUserDao().Find(r.Context(), dao.Where{})
	if err != nil {
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}
	renders := []render.Renderer{}
	for _, user := range users {
		renders = append(renders, user)
	}
	render.Status(r, http.StatusOK)
	if err := render.Render(rw, r, &res.JSON{"users": renders}); err != nil {
		render.Render(rw, r, res.ErrRender(err))
		return
	}
}

func (h *UserHandler) FindById(rw http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userIDCtxKey{}).(string)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}
	user, err := h.daoFactory.GetUserDao().FindOne(r.Context(), dao.Where{{Key: "_id", Value: oid}})
	if err != nil {
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	if err := render.Render(rw, r, &res.JSON{"user": user}); err != nil {
		render.Render(rw, r, res.ErrRender(err))
		return
	}
}

func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), userIDCtxKey{}, chi.URLParam(r, "id"))
		next.ServeHTTP(w, r.Clone(ctx))
	})
}

func NewUserHandler(logger *logger.Logger, factory *dao.Factory) *UserHandler {
	return &UserHandler{
		logger:     logger,
		daoFactory: factory,
	}
}
