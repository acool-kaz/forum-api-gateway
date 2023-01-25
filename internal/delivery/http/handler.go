package http

import (
	"log"
	"net/http"

	"github.com/acool-kaz/forum-api-gateway/internal/config"
	auth_svc "github.com/acool-kaz/forum-api-gateway/internal/delivery/http/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	authService *auth_svc.AuthService
}

func InitHandler(cfg *config.Config) (*Handler, error) {
	log.Println("init http handler")

	authService, err := auth_svc.InitAuthService(cfg)
	if err != nil {
		return nil, err
	}

	return &Handler{
		authService: authService,
	}, nil
}

func (h *Handler) InitRoutes() http.Handler {
	log.Println("init routes")

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	h.authService.RegisterAuthServiceRoutes(router)

	return router
}
