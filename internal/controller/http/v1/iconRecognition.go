package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
)

type iconRecognitionRoutesRoutes struct {
	t usecase.Translation
	l logger.Interface
}

func newIconRecognitionRoutes(handler *gin.RouterGroup, t usecase.Translation, l logger.Interface) {
	r := &iconRecognitionRoutesRoutes{t, l}

	h := handler.Group("")
	{
		h.GET("/doRecognition", r.uploadDashboard)
		h.POST("/doRecognition", r.doRecognition)
	}
}

type doRecognitionResponse struct {
	recognitionResult []string
}

func (r *iconRecognitionRoutesRoutes) doRecognition(c *gin.Context) {
	// Получим файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка загрузки файла: %s", err.Error()))
		return
	}

	// сохраняем файл на сервере
	err = c.SaveUploadedFile(file, "./pkg/file_storage/test/"+file.Filename)

	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Ошибка сохранения файла: %s", err.Error()))
		return
	}
	//TODO обращение к usecase

	//TODO вернуть распознанное

}

func (r *iconRecognitionRoutesRoutes) uploadDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "mainScreen.html", gin.H{
		"block_title": "Test page",
	})
}
