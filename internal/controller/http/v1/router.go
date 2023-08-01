// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/evrone/go-clean-template/internal/controller/http/middleware"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	au *middleware.Auth
}

func NewRouter(au *middleware.Auth) *Router {
	return &Router{
		au: au,
	}
}

// InitRoutes -.
func (r *Router) InitRoutes(handler *gin.Engine, l logger.Interface, t usecase.Recognition) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Static файлы

	handler.Static("/internal/controller/static/templates/css", "./internal/controller/static/templates/css")
	handler.Static("/pkg/file_storage/detected", "./pkg/file_storage/detected")
	handler.LoadHTMLGlob("./internal/controller/static/templates/html/*.html")

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	router := handler.Group("/v1")
	routerPath := "/v1"
	{
		subGroup1URL := "/private"
		private := router.Group(subGroup1URL)
		private.Use(r.au.AuthRequired())
		nextURL := routerPath + subGroup1URL + "/api/doRecognition"
		newIconRecognitionRoutes(private, t, l)

		subGroup2URL := "/auth"
		auth := router.Group(subGroup2URL)
		authPath := routerPath + subGroup2URL
		r.au.NewAuthRoutes(auth, authPath, nextURL)
	}
}
