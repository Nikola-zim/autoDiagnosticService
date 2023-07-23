package v1

import (
	"fmt"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
)

type recognition struct {
	useCase usecase.Recognition
	l       logger.Interface
}

func newIconRecognitionRoutes(handler *gin.RouterGroup, t usecase.Recognition, l logger.Interface) {
	r := &recognition{t, l}

	h := handler.Group("/api")
	{
		h.GET("/doRecognition", r.uploadDashboard)
		h.POST("/doRecognition", r.doRecognition)
	}

	u := handler.Group("/user")
	{
		u.GET("/register", r.authPage)
		u.POST("/register", r.register)
		u.GET("/login", r.authPage)
		u.POST("/login", r.login)
	}

	b := handler.Group("/balance")
	{
		b.GET("/sum", r.uploadDashboard)
		b.POST("/add", r.doRecognition)
	}
}

type doRecognitionResponse struct {
	recognitionResult []string
}

func (r *recognition) doRecognition(c *gin.Context) {
	// Получим файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка загрузки файла: %s", err.Error()))
		return
	}

	// сохраняем файл на сервере
	err = c.SaveUploadedFile(file, "./pkg/file_storage/images/"+file.Filename)

	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Ошибка сохранения файла: %s", err.Error()))
		return
	}

}

func (r *recognition) uploadDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "mainScreen.html", gin.H{
		"block_title": "Test page",
	})
}

func (r *recognition) authPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth.html", gin.H{
		"block_title": "Test page",
	})
}

func (r *recognition) register(c *gin.Context) {
	c.HTML(http.StatusOK, "auth.html", gin.H{
		"block_title": "Test page",
	})
	var user entity.User
	user.Login = c.PostForm("uname")
	user.Password = c.PostForm("psw")

	r.l.Info("Username: %s; Password: %s", user.Login, user.Password)

	err := r.useCase.AddUser(c, user)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Wrong input or user already exists: %s", err.Error()))
		return
	}
}

func (r *recognition) login(c *gin.Context) {
	c.HTML(http.StatusOK, "auth.html", gin.H{
		"block_title": "Test page",
	})
}
