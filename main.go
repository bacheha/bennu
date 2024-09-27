package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/knuls/bennu/app"
	"github.com/knuls/bennu/dao"
	"github.com/knuls/bennu/handlers"
	"github.com/knuls/horus/config"
	"github.com/knuls/horus/logger"
	"github.com/knuls/horus/middlewares"
	"github.com/knuls/horus/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// logger
	log, err := logger.NewZapLogger()
	if err != nil {
		fmt.Printf("logger new error: %v", err)
		os.Exit(1)
	}
	defer log.GetZapLogger().Sync()

	// config
	c := config.New("bennu")
	c.SetBindings(app.Bindings)
	var cfg *app.Config
	if err := c.Load("config", "yaml", ".", &cfg); err != nil {
		log.Error("config load", "error", err)
		return
	}

	// db
	dbCtx, cancel := context.WithTimeout(context.Background(), cfg.Store.Timeout*time.Second)
	defer cancel()
	uri := fmt.Sprintf("%s://%s:%d", cfg.Store.Client, cfg.Store.Host, cfg.Store.Port)
	client, err := mongo.Connect(dbCtx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Error("db connect", "error", err)
		return
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Error("db disconnect", "error", err)
			return
		}
	}()
	pingCtx, cancel := context.WithTimeout(context.Background(), cfg.Store.Timeout*time.Second)
	defer cancel()
	if err = client.Ping(pingCtx, readpref.Primary()); err != nil {
		log.Error("db ping", "error", err)
		return
	}

	// validator
	v, err := validator.New()
	if err != nil {
		log.Error("validator new", "error", err)
		return
	}

	// dao factory
	db := client.Database(cfg.Store.Name)
	factory := dao.NewDaoFactory(db, v)

	// mux
	mux := chi.NewRouter()

	// middlewares
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Security.Allowed.Origins,
		AllowedMethods:   cfg.Security.Allowed.Methods,
		AllowedHeaders:   cfg.Security.Allowed.Headers,
		AllowCredentials: cfg.Security.AllowCredentials,
	}))
	mux.Use(middlewares.JSON)
	mux.Use(middlewares.RealIP)
	mux.Use(middlewares.RequestID)
	mux.Use(middlewares.Recoverer)
	mux.Use(middlewares.Logger(log))

	// handlers
	mux.Mount("/user", handlers.NewUserHandler(log, factory).Routes())
	mux.Mount("/organization", handlers.NewOrganizationHandler(log, factory).Routes())
	mux.Mount("/auth", handlers.NewAuthHandler(log, factory, cfg).Routes())

	// server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Service.Port),
		Handler:      mux,
		ErrorLog:     log.GetStdLogger(),
		ReadTimeout:  cfg.Server.Timeout.Read * time.Second,
		WriteTimeout: cfg.Server.Timeout.Write * time.Second,
		IdleTimeout:  cfg.Server.Timeout.Idle * time.Second,
	}

	// listen
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("listen and serve", "error", err)
			return
		}
	}()
	log.Infof("starting %s service on port: %d", cfg.Service.Name, cfg.Service.Port)

	// shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	sig := <-sigCh
	log.Infof("signal: %s", sig.String())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.Timeout.Shutdown*time.Second)
	defer cancel()
	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Error("shutdown", "error", err)
		return
	}
}
