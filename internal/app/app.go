// Package app configures and runs application.
package app

import (
	"autoDiagnosticService/internal/controller/http/middlewares"
	"autoDiagnosticService/internal/entity"
	"autoDiagnosticService/internal/file_storage"
	"autoDiagnosticService/internal/usecase/repo"
	"autoDiagnosticService/internal/usecase/worker"
	"autoDiagnosticService/pkg/postgres"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"autoDiagnosticService/config"
	v1 "autoDiagnosticService/internal/controller/http/v1"
	"autoDiagnosticService/internal/usecase"
	"autoDiagnosticService/pkg/httpserver"
	"autoDiagnosticService/pkg/logger"
	_ "autoDiagnosticService/pkg/postgres"

	"autoDiagnosticService/internal/controller/telegram"
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

	// file_storage
	fs := fileStorage.New(cfg.Storage.ImagePath)

	// Detection Worker
	detectionWorker := worker.NewDetectionWebAPI(detectionUseCase, fs, newAnswer, cfg.Detector)
	go func() {
		err = detectionWorker.Run(context.Background())
		if err != nil {
			l.Fatal(err)
		}
	}()
	//
	classes := entity.NewClasses()
	//telegram bot
	telegramBot, err := telegram.New(cfg.TG.BotToken, cfg.TG.ImagePath, detectionUseCase, classes, newAnswer)
	if err != nil {
		l.Warn(fmt.Sprintf("app - Run - telegram.New: %w", err))
	}

	// HTTP Server
	handler := gin.New()
	au := middlewares.NewAuth(l)
	router := v1.NewRouter(au)
	router.InitRoutes(handler, l, detectionUseCase, cfg.TG.ImagePath)
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
