package post

import (
	"log"
	"net/http"
	"strconv"

	"github.com/acool-kaz/forum-api-gateway/internal/config"
	"github.com/acool-kaz/forum-api-gateway/internal/delivery/http/auth"
	"github.com/acool-kaz/forum-api-gateway/internal/models"
	"github.com/acool-kaz/forum-api-gateway/pkg/json"
	"github.com/go-chi/chi/v5"
)

type PostService struct {
	authService *auth.AuthService
}

func InitPostService(cfg *config.Config, authService *auth.AuthService) (*PostService, error) {
	return &PostService{
		authService: authService,
	}, nil
}

func (p *PostService) RegisterPostServiceRoutes(router *chi.Mux) {
	log.Println("register post service endpoints")

	router.Route("/api/post", func(post chi.Router) {
		post.Use(p.authService.AuthMiddleware)

		post.Get("/", func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(models.CurrentUser)

			json.SendJson(w, map[string]interface{}{
				"hello": "world from " + strconv.Itoa((int(userId.(uint)))),
			})
		})
	})
}
