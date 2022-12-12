package handlers

import (
	"github.com/bachehah/horus/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrganizationHandler struct {
	Logger   *logger.Logger
	Validate *validator.Validate
	Client   *mongo.Client
}

func (h *OrganizationHandler) Routes() *chi.Mux {
	mux := chi.NewRouter()
	mux.Get("/", nil)  // GET /organization
	mux.Post("/", nil) // POST /organization
	mux.Route("/{id}", func(mux chi.Router) {
		mux.Get("/", nil) // GET /organization/:id
	})
	return mux
}

func NewOrganizationHandler(logger *logger.Logger, validate *validator.Validate, client *mongo.Client) *UserHandler {
	return &UserHandler{
		Logger:   logger,
		Validate: validate,
		Client:   client,
	}
}
