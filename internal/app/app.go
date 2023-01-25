package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/acool-kaz/forum-api-gateway/internal/config"

	httpHandler "github.com/acool-kaz/forum-api-gateway/internal/delivery/http"
)

type app struct {
	cfg *config.Config

	httpServer  *http.Server
	httpHandler *httpHandler.Handler
}

func InitApp(cfg *config.Config) (*app, error) {
	log.Println("init app")

	httpHandler, err := httpHandler.InitHandler(cfg)
	if err != nil {
		return nil, err
	}

	return &app{
		cfg:         cfg,
		httpHandler: httpHandler,
	}, nil
}

func (a *app) Run() {
	log.Println("run app")

	go func() {
		if err := a.startHTTP(); err != nil {
			log.Println(err)
			return
		}
	}()
	log.Println("http server started on", a.cfg.Http.Port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	fmt.Println()
	log.Println("Received terminate, graceful shutdown", sig)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Println(err)
		return
	}
}

func (a *app) startHTTP() error {
	handler := a.httpHandler.InitRoutes()

	a.httpServer = &http.Server{
		Addr:         ":" + a.cfg.Http.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * time.Duration(a.cfg.Http.Read),
		WriteTimeout: time.Second * time.Duration(a.cfg.Http.Write),
	}

	return a.httpServer.ListenAndServe()
}
