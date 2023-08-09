package worker

import (
	"autoDiagnosticService/config"
	"autoDiagnosticService/internal/entity"
	"autoDiagnosticService/internal/usecase"
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

// DetectionWebAPI -.
type DetectionWebAPI struct {
	scheduler *gocron.Scheduler
	useCase   usecase.Recognition
	newAnswer chan bool
	config    config.Detector
}

const (
	defaultLocation          = "Europe/Moscow"
	defaultDetectedImagePath = "internal/file_storage/detected"
	layout                   = "2006_01_02"
)

// NewDetectionWebAPI -.
func NewDetectionWebAPI(useCase usecase.Recognition, newAnswer chan bool, config config.Detector) *DetectionWebAPI {
	location, _ := time.LoadLocation(defaultLocation)
	return &DetectionWebAPI{
		useCase:   useCase,
		scheduler: gocron.NewScheduler(location),
		newAnswer: newAnswer,
		config:    config,
	}
}

func (dw *DetectionWebAPI) Run(ctx context.Context) error {
	recognitionTask := func(ctx context.Context) error {
		tasks, err := dw.useCase.GetRecognitionTasks(ctx)
		if err != nil {
			log.Info().Msgf("no tasks or error", err)
			return err
		}
		err = dw.serverRecognitionQuery(ctx, tasks)
		if err != nil {
			log.Info().Msgf("failed to get results from recognition server", err)
			return err
		}
		// Send answers
		dw.newAnswer <- true
		return nil
	}
	_, err := dw.scheduler.Every(2).Second().Do(recognitionTask, ctx)
	if err != nil {
		return err
	}
	dw.scheduler.StartAsync()

	<-ctx.Done()

	return nil
}

func (dw *DetectionWebAPI) serverRecognitionQuery(ctx context.Context, tasks []entity.Request) error {
	for _, task := range tasks {
		// Открываем файл с изображением
		file, err := os.Open(task.ImagePathName)
		if err != nil {
			log.Info().Msgf("with %s%s ", task.ImagePathName, err)
			continue
		}
		defer file.Close()

		// Создаем буфер для формы
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		// Добавляем файл в форму
		part, err := writer.CreateFormFile(dw.config.FormField, dw.config.FormName)
		if err != nil {
			return err
		}
		_, err = io.Copy(part, file)

		// Закрываем форму
		writer.Close()

		// Создаем запрос POST
		req, err := http.NewRequest("POST", dw.config.URL, &buf)
		if err != nil {
			return err
		}

		// Устанавливаем заголовок Content-Type для формы
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Отправляем запрос и получаем ответ
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Читаем ответ
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			panic(err)
		}
		var t entity.DetectionResult
		err = json.Unmarshal(body, &t)
		task.Description = t.Description

		if err != nil {
			return err
		}

		// Save image
		newPath, err := imageSave(t.Image, task.ImagePathName)
		if err != nil {
			return err
		}
		task.ResImgPathName = newPath
		// Change status
		err = dw.useCase.MakeRecognized(ctx, task)
		if err != nil {
			return err
		}

	}
	return nil
}

func imageSave(img string, pathName string) (string, error) {
	date := time.Now().Format(layout)
	decodedBytes, err := b64.StdEncoding.DecodeString(img)

	if err != nil {
		return "", err
	}

	// Создаем файл на сервере
	currentPath := defaultDetectedImagePath + "/" + date
	if _, err := os.Stat(currentPath); os.IsNotExist(err) {
		os.MkdirAll(currentPath, 0700) // Create your file
	}
	// Save image
	ss := strings.Split(pathName, "/")
	s := ss[len(ss)-1]
	filePathAndName := currentPath + "/" + s
	out, err := os.Create(filePathAndName)
	if err != nil {
		return "", err
	}
	defer out.Close()
	// Сохраняем декодированные байты в файл
	err = os.WriteFile(filePathAndName, decodedBytes, 0644)
	if err != nil {
		log.Err(err).Msgf("ioutil.WriteFile fails on: %s", filePathAndName)
	}

	return filePathAndName, nil
}
