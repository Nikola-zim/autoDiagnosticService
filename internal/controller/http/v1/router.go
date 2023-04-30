// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewRouter -.
func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.Translation) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Static файлы
	handler.Static("/internal/controller/static/templates/css", "./internal/controller/static/templates/css")
	handler.LoadHTMLGlob("./internal/controller/static/templates/html/*.html")

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		newTranslationRoutes(h, t, l)
	}
	iconRecognition := handler.Group("/iconRecognition")
	{
		newIconRecognitionRoutes(iconRecognition, t, l)
	}
}
