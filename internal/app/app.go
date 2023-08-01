// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"github.com/evrone/go-clean-template/internal/controller/http/middleware"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase/repo"
	"github.com/evrone/go-clean-template/internal/usecase/worker"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/evrone/go-clean-template/config"
	v1 "github.com/evrone/go-clean-template/internal/controller/http/v1"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/logger"
	_ "github.com/evrone/go-clean-template/pkg/postgres"

	"github.com/evrone/go-clean-template/internal/controller/telegram"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	newAnswer := make(chan bool)

	// Use case
	detectionUseCase := usecase.New(
		repo.NewRecognitionRepo(pg),
		repo.NewAuth(pg),
	)

	// Detection Worker
	detectionWorker := worker.NewDetectionWebAPI(detectionUseCase, newAnswer)
	go func() {
		err = detectionWorker.Run(context.Background())
		if err != nil {
			l.Fatal(err)
		}
	}()
	//
	classes := entity.NewClasses()
	//telegram bot
	telegramBot, err := telegram.New(cfg.TG.BotToken, detectionUseCase, classes, newAnswer)
	if err != nil {
		l.Warn(fmt.Sprintf("app - Run - telegram.New: %w", err))
	}

	// HTTP Server
	handler := gin.New()
	// Setup the cookie store for session management
	var secret = []byte("secret")
	handler.Use(sessions.Sessions("mysession", cookie.NewStore(secret)))
	au := middleware.NewAuth(l, detectionUseCase)
	router := v1.NewRouter(au)
	router.InitRoutes(handler, l, detectionUseCase)
	httpServer := httpserver.New(handler, telegramBot, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
