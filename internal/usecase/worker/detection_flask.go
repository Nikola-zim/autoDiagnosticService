package worker

import (
	"autoDiagnosticService/config"
	"autoDiagnosticService/internal/entity"
	"autoDiagnosticService/internal/file_storage"
	"autoDiagnosticService/internal/usecase"
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// DetectionWebAPI -.
type DetectionWebAPI struct {
	scheduler   *gocron.Scheduler
	useCase     usecase.Recognition
	fileStorage *fileStorage.FileStorage
	newAnswer   chan bool
	config      config.Detector
}

const (
	defaultLocation = "Europe/Moscow"
)

// NewDetectionWebAPI -.
func NewDetectionWebAPI(useCase usecase.Recognition, fileStorage *fileStorage.FileStorage, newAnswer chan bool, config config.Detector) *DetectionWebAPI {
	location, _ := time.LoadLocation(defaultLocation)
	return &DetectionWebAPI{
		useCase:     useCase,
		fileStorage: fileStorage,
		scheduler:   gocron.NewScheduler(location),
		newAnswer:   newAnswer,
		config:      config,
	}
}

func (dw *DetectionWebAPI) Run(ctx context.Context) error {
	recognitionTask := func(ctx context.Context) error {
		tasks, err := dw.useCase.GetRecognitionTasks(ctx)
		if err != nil {
			log.Info().Msgf("no tasks or error: %s", err)
			return err
		}
		err = dw.serverRecognitionQuery(ctx, tasks)
		if err != nil {
			log.Info().Msgf("failed to get results from recognition server: %s", err)
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
			log.Info().Msgf("run with error; image: %s; error: %s ", task.ImagePathName, err)
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
		newPath, err := dw.fileStorage.ImageSave(t.Image, task.ImagePathName)
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
