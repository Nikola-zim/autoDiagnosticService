package v1

import (
	"autoDiagnosticService/internal/entity"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"image"
	"log"
	"net/http"
	"os"
	"time"

	"autoDiagnosticService/pkg/logger"
	"github.com/disintegration/imaging"
)

const (
	layout   = "2006_01_02"
	filePath = "./internal/file_storage/images/"
)

type recognition struct {
	useCase Recognition
	l       logger.Interface
}

func newIconRecognitionRoutes(handler *gin.RouterGroup, t Recognition, l logger.Interface) {
	r := &recognition{t, l}

	h := handler.Group("/api")
	{
		h.GET("/recognized_images", r.recognizedImages)
		h.POST("/recognize", r.uploadImage)
	}

	b := handler.Group("/balance")
	{
		b.POST("/add", r.addPoints)
	}
}

type doRecognitionResponse struct {
	recognitionResult []string
}

func (r *recognition) uploadImage(c *gin.Context) {
	// Получим файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка загрузки файла: %s", err.Error()))
		return
	}

	// сохраняем файл на сервере
	date := time.Now().Format(layout)
	filePathName := fmt.Sprintf(filePath+"%s"+"/%s", date, file.Filename)
	// Создаем dir
	if _, err := os.Stat(filePath + "/" + date); os.IsNotExist(err) {
		os.MkdirAll(filePath+"/"+date, 0700) // Create your file
	}
	err = c.SaveUploadedFile(file, filePathName)

	if err != nil {
		c.String(http.StatusBadRequest, "an error occurred while saving the file: %v", err)
		return
	}
	// Open the file.
	f, err := os.Open(filePathName)
	if err != nil {
		c.String(http.StatusBadRequest, "an error occurred while opening the file: %v", err)
		return
	}
	img, _, err := image.Decode(f)
	if err != nil {
		c.String(http.StatusBadRequest, "an error occurred while decoding the image: %v", err)
		return
	}

	// Resize the image.
	resizedImg := imaging.Resize(img, 640, 640, imaging.Lanczos)

	// Save the resized image.
	err = imaging.Save(resizedImg, filePathName)
	if err != nil {
		c.String(http.StatusBadRequest, "an error occurred while saving the resized image: %v", err)
		return
	}
	// Add record to pg for worker
	session := sessions.Default(c)
	user := session.Get(userKey)
	newReq := entity.Request{
		UserID:        fmt.Sprintf("%v", user),
		ImagePathName: filePathName,
	}
	err = r.useCase.AddRequest(c, newReq)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, fmt.Sprintf("%s", err))
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Image uploaded!"})
	}

}

func (r *recognition) recognizedImages(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	results, err := r.useCase.GetRecognitionAnswersWEB(c, fmt.Sprintf("%v", user))
	if err != nil {
		log.Println(err)
	}

	images := make([]string, 0, len(results))
	for _, v := range results {
		images = append(images, v.ResImgPathName)
	}

	c.HTML(http.StatusOK, "recognition.html", gin.H{
		"block_title": "Test page",
		"Images":      images,
	})
}

func (r *recognition) addPoints(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	var balanceAdd entity.Balance
	err := c.ShouldBindJSON(&balanceAdd)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = r.useCase.AddPoints(c, balanceAdd.Points, fmt.Sprintf("%v", user))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s", err))
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance is replenished!"})
}
