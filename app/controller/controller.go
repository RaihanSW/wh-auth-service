package controller

import (
	"dc-transaction-service/app/controller/handler"
	"dc-transaction-service/app/controller/handler/middleware"
	"dc-transaction-service/app/model/constant"
	"dc-transaction-service/config"
	"dc-transaction-service/docs"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Controller struct {
	generalHandler handler.GeneralHandlerInterface
	authHandler    handler.AuthHandlerInterface
	roleHandler    handler.RoleHandlerInterface
	cfg            *config.ApplicationConfig
	env            string
}

func NewController(
	generalHandler handler.GeneralHandlerInterface,
	authHandler handler.AuthHandlerInterface,
	roleHandler handler.RoleHandlerInterface,
	cfg *config.ApplicationConfig,
	env string) *Controller {

	return &Controller{
		generalHandler: generalHandler,
		authHandler:    authHandler,
		roleHandler:    roleHandler,
		cfg:            cfg,
		env:            env,
	}
}

func (c *Controller) StartRoute() *http.Server {
	r := chi.NewRouter()

	r.Use(middleware.Middleware())

	r.Get("/ping", middleware.SetRequestID(c.generalHandler.Ping))
	r.Get("/health/check", middleware.SetRequestID(c.generalHandler.HealthCheck))
	r.Get("/health/check-db", middleware.SetRequestID(c.generalHandler.HealthCheckDB))

	// Init Swagger with swaggo
	docs.SwaggerInfo.Title = "Warehouse Auth Service"
	docs.SwaggerInfo.Description = fmt.Sprintf("Warehouse Auth Service with Golang - Environtment: %s", strings.ToLower(c.env))
	docs.SwaggerInfo.Version = "1.1"
	switch strings.ToLower(c.env) {
	case constant.EnvLocal:
		docs.SwaggerInfo.Schemes = []string{"http"}
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", c.cfg.Port.Service)
		docs.SwaggerInfo.BasePath = "/"
	case constant.EnvDev:
		docs.SwaggerInfo.Schemes = []string{"https"}
		docs.SwaggerInfo.Host = ""
		docs.SwaggerInfo.BasePath = "/transaction/"
	case constant.EnvStaging:
		docs.SwaggerInfo.Schemes = []string{"https"}
		docs.SwaggerInfo.Host = ""
		docs.SwaggerInfo.BasePath = "/transaction/"
	}

	if strings.ToLower(c.env) != constant.EnvProd {
		r.Mount("/swagger", httpSwagger.WrapHandler)
	}
	// Account routes (register / login / logout + user CRUD — table "User".users)
	r.Route("/account", func(r chi.Router) {
		r.Use(middleware.GetLogMiddleware(*c.cfg))
		r.Post("/register", c.authHandler.AuthRegister)
		r.Post("/login", c.authHandler.AuthLogin)
		r.Post("/logout", c.authHandler.AuthLogout)
		r.Get("/{uid}", c.authHandler.UserGetByID)
		r.Put("/{uid}", c.authHandler.UserUpdate)
		r.Delete("/{uid}", c.authHandler.UserDelete)
	})

	// Role routes (table "User".roles)
	r.Route("/role", func(r chi.Router) {
		r.Use(middleware.GetLogMiddleware(*c.cfg))
		r.Post("/", c.roleHandler.RoleCreate)
		r.Get("/", c.roleHandler.RoleGetAll)
		r.Get("/{id}", c.roleHandler.RoleGetByID)
		r.Put("/{id}", c.roleHandler.RoleUpdate)
		r.Delete("/{id}", c.roleHandler.RoleDelete)
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", c.cfg.Port.Service),
		WriteTimeout: c.cfg.Port.ServiceTimeout * time.Second,
		ReadTimeout:  c.cfg.Port.ServiceTimeout * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return srv
}
