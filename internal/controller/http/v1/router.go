// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"autoDiagnosticService/internal/controller/http/middlewares"
	"autoDiagnosticService/internal/usecase"
	"autoDiagnosticService/pkg/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var secret = []byte("secret")

type Router struct {
	au *middlewares.Auth
}

func NewRouter(au *middlewares.Auth) *Router {
	return &Router{
		au: au,
	}
}

// InitRoutes -.
func (r *Router) InitRoutes(handler *gin.Engine, l logger.Interface, t usecase.Recognition, storagePath string) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.Use(sessions.Sessions("mysession", cookie.NewStore(secret)))

	// Static файлы

	handler.Static("/internal/controller/static/templates/css", "./internal/controller/static/templates/css")
	handler.Static("/internal/file_storage/detected", "./internal/file_storage/detected")
	handler.LoadHTMLGlob("./internal/controller/static/templates/html/*.html")

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	router := handler.Group("/v1")
	{
		subGroup1URL := "/private"
		private := router.Group(subGroup1URL)
		private.Use(r.au.AuthRequired())
		newIconRecognitionRoutes(private, t, l, storagePath)

		subGroup2URL := "/auth"
		auth := router.Group(subGroup2URL)
		NewAuthHandlers(auth, t, l)
	}
}
