package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/knuls/bennu/dao"
	"github.com/knuls/bennu/organizations"
	"github.com/knuls/horus/logger"
	"github.com/knuls/horus/middlewares"
	"github.com/knuls/horus/res"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type organizationIDCtxKey struct{}

type organizationHandler struct {
	logger     logger.Logger
	daoFactory dao.Factory
}

func NewOrganizationHandler(logger logger.Logger, factory dao.Factory) *organizationHandler {
	return &organizationHandler{
		logger:     logger,
		daoFactory: factory,
	}
}

func (h *organizationHandler) Routes() *chi.Mux {
	mux := chi.NewRouter()
	mux.Get("/", h.Find)    // GET /organization
	mux.Post("/", h.Create) // POST /organization
	mux.Route("/{id}", func(mux chi.Router) {
		mux.Use(middlewares.ValidateObjectID("id"))
		mux.Use(OrganizationCtx)
		mux.Get("/", h.FindById) // GET /organization/:id
	})
	return mux
}

func (h *organizationHandler) Find(rw http.ResponseWriter, r *http.Request) {
	orgs, err := h.daoFactory.GetOrganizationDao().Find(r.Context(), dao.Where{})
	if err != nil {
		h.logger.Error("failed to find organizations", "error", err)
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}
	renders := []render.Renderer{}
	for _, org := range orgs {
		renders = append(renders, org)
	}
	render.Status(r, http.StatusOK)
	if err = render.Render(rw, r, &res.JSON{"organizations": renders}); err != nil {
		h.logger.Error("failed to render", "error", err)
		render.Render(rw, r, res.ErrRender(err))
		return
	}
}

func (h *organizationHandler) Create(rw http.ResponseWriter, r *http.Request) {
	org := organizations.NewOrganization()
	defer r.Body.Close()
	if err := org.FromJSON(r.Body); err != nil {
		h.logger.Error("failed to decode request body", "error", err)
		render.Render(rw, r, res.ErrDecode(err))
		return
	}
	id, err := h.daoFactory.GetOrganizationDao().Create(r.Context(), org)
	if err != nil {
		h.logger.Error("failed to create organization", "error", err)
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
	if err = render.Render(rw, r, &res.JSON{"id": id}); err != nil {
		h.logger.Error("failed to render", "error", err)
	}
}

func (h *organizationHandler) FindById(rw http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(organizationIDCtxKey{}).(string)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.logger.Error("failed to convert hex to object id", "error", err)
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}
	org, err := h.daoFactory.GetOrganizationDao().FindOne(r.Context(), dao.Where{{Key: "_id", Value: oid}})
	if err != nil {
		h.logger.Error("failed to find organization", "error", err)
		render.Render(rw, r, res.ErrBadRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	if err := render.Render(rw, r, &res.JSON{"organization": org}); err != nil {
		h.logger.Error("failed to render", "error", err)
		render.Render(rw, r, res.ErrRender(err))
		return
	}
}

func OrganizationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), organizationIDCtxKey{}, chi.URLParam(r, "id"))
		next.ServeHTTP(w, r.Clone(ctx))
	})
}
