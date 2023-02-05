package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/acool-kaz/forum-api-gateway/internal/config"
	"github.com/acool-kaz/forum-api-gateway/internal/models"
	auth_svc_pb "github.com/acool-kaz/forum-api-gateway/pkg/auth_svc/pb"
	"github.com/acool-kaz/forum-api-gateway/pkg/json"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthService struct {
	Client auth_svc_pb.AuthServiceClient
}

func InitAuthService(cfg *config.Config) (*AuthService, error) {
	log.Println("init auth service client")

	conn, err := grpc.Dial(cfg.AuthService.Host+":"+cfg.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &AuthService{
		Client: auth_svc_pb.NewAuthServiceClient(conn),
	}, nil
}

func (a *AuthService) RegisterAuthServiceRoutes(router *chi.Mux) {
	log.Println("register auth service endpoints")

	router.Route("/auth", func(auth chi.Router) {
		auth.Post("/sign-up", a.signUpHandler)
		auth.Post("/sign-in", a.signInHandler)
	})
}

func (a *AuthService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header, ok := r.Header["Authorization"]
		if !ok {
			json.SendError(w, errors.New("empty auth header"))
			return
		}

		headerParts := strings.Split(header[0], " ")
		if len(headerParts) != 2 {
			json.SendError(w, errors.New("invalid auth header"))
			return
		}

		resp, err := a.Client.Validate(r.Context(), &auth_svc_pb.ValidateRequest{
			AccessToken: headerParts[1],
		})
		if err != nil {
			json.SendError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), models.CurrentUser, uint(resp.GetUserId()))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthService) signUpHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserSignUp
	err := json.ParseJson(r, &user)
	if err != nil {
		json.SendError(w, err)
		return
	}

	req := auth_svc_pb.RegisterRequest{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		Password:  user.Password,
	}

	resp, err := a.Client.Register(r.Context(), &req)
	if err != nil {
		json.SendError(w, err)
		return
	}

	err = json.SendJson(w, resp)
	if err != nil {
		json.SendError(w, err)
		return
	}
}

func (a *AuthService) signInHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserSignIn
	err := json.ParseJson(r, &user)
	if err != nil {
		json.SendError(w, err)
		return
	}

	req := auth_svc_pb.LoginRequest{
		Email:    &user.Email,
		Username: &user.Username,
		Password: user.Password,
	}

	resp, err := a.Client.Login(r.Context(), &req)
	if err != nil {
		json.SendError(w, err)
		return
	}

	err = json.SendJson(w, resp)
	if err != nil {
		json.SendError(w, err)
		return
	}
}
