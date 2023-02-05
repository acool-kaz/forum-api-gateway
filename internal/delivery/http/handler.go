package http

import (
	"log"
	"net/http"

	"github.com/acool-kaz/forum-api-gateway/internal/config"
	auth_svc "github.com/acool-kaz/forum-api-gateway/internal/delivery/http/auth"
	post_svc "github.com/acool-kaz/forum-api-gateway/internal/delivery/http/post"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	authService *auth_svc.AuthService
	postService *post_svc.PostService
}

func InitHandler(cfg *config.Config) (*Handler, error) {
	log.Println("init http handler")

	authService, err := auth_svc.InitAuthService(cfg)
	if err != nil {
		return nil, err
	}

	postService, err := post_svc.InitPostService(cfg, authService)
	if err != nil {
		return nil, err
	}

	return &Handler{
		authService: authService,
		postService: postService,
	}, nil
}

func (h *Handler) InitRoutes() http.Handler {
	log.Println("init routes")

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	h.authService.RegisterAuthServiceRoutes(router)
	h.postService.RegisterPostServiceRoutes(router)

	return router
}
