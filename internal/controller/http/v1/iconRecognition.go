package v1

import (
	"fmt"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
)

const (
	layout   = "2006_01_02"
	filePath = "./internal/file_storage/images/"
)

type recognition struct {
	useCase usecase.Recognition
	l       logger.Interface
}

func newIconRecognitionRoutes(handler *gin.RouterGroup, t usecase.Recognition, l logger.Interface) {
	r := &recognition{t, l}

	h := handler.Group("/api")
	{
		h.GET("/recognized_images", r.recognizedImages)
		h.POST("/recognize", r.uploadImage)
	}

	b := handler.Group("/balance")
	{
		b.GET("/sum", r.pointsSum)
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
		log.Println(err)
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

func (r *recognition) addPoints(context *gin.Context) {

}

func (r *recognition) pointsSum(context *gin.Context) {

}

// Функция для загрузки файла с указанного URL
func downloadFile(url string, messageID string) (string, error) {
	date := time.Now().Format(layout)
	fileName := date + "/to_detect" + messageID + ".jpg"
	filePathAndName := fmt.Sprintf(filePath+"%s", fileName)
	// Создаем dir
	if _, err := os.Stat(filePath + "/" + date); os.IsNotExist(err) {
		os.MkdirAll(filePath+"/"+date, 0700) // Create your file
	}

	// Save image
	out, err := os.Create(filePathAndName)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Получаем содержимое файла по URL и записываем в созданный файл
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	receivedImage, _, err := image.Decode(resp.Body)
	// Resize to 640x640
	newImage := resize.Resize(640, 640, receivedImage, resize.Lanczos3)
	// Save image
	err = jpeg.Encode(out, newImage, nil)
	if err != nil {
		return "", err
	}
	return filePathAndName, nil
}
