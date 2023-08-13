package telegram

import (
	"autoDiagnosticService/internal/entity"
	"autoDiagnosticService/internal/usecase"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Telegram struct {
	BotAPI    *tgbotapi.BotAPI
	botToken  string
	useCase   usecase.Recognition
	newAnswer chan bool
	classes   *entity.Classes
	filePath  string
}

const layout = "2006_01_02"

func New(token string, filePath string, useCase usecase.Recognition, classes *entity.Classes, newAnswer chan bool, opts ...Option) (*Telegram, error) {

	botAPI, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, err
	}

	tg := &Telegram{
		botToken:  token,
		BotAPI:    botAPI,
		useCase:   useCase,
		newAnswer: newAnswer,
		classes:   classes,
		filePath:  filePath,
	}
	// Custom options
	for _, opt := range opts {
		opt(tg)
	}
	log.Printf("Authorized on account %s", tg.BotAPI.Self.UserName)

	return tg, err
}

func (tg *Telegram) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tg.BotAPI.GetUpdatesChan(u)
	ctx := context.Background()
	go func() {
		for update := range updates {

			if update.Message != nil { // If we got a message
				// Если сообщение содержит изображение
				if update.Message.Photo != nil {
					// Получаем информацию об изображении
					fileID := (update.Message.Photo)[len(update.Message.Photo)-1].FileID
					file, err := tg.BotAPI.GetFile(tgbotapi.FileConfig{FileID: fileID})
					if err != nil {
						log.Println(err)
						continue
					}

					// Загружаем изображение
					imagePathName, err := tg.downloadFile(file.Link(tg.BotAPI.Token), strconv.Itoa(update.Message.MessageID))
					if err != nil {
						log.Println(err)
						continue
					}

					// Add record to pg for worker
					newReq := entity.Request{
						ChatID:        update.Message.Chat.ID,
						ImagePathName: imagePathName,
					}
					err = tg.useCase.AddRequest(ctx, newReq)
					if err != nil {
						log.Println(err)
					}
					// Отправляем сообщение об успешной загрузке изображения
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Изображение успешно загружено на сервер....")
					_, err = tg.BotAPI.Send(msg)
					if err != nil {
						log.Println(err)
					}

				}
			}
		}
	}()
	for {
		if <-tg.newAnswer {
			results, err := tg.useCase.GetRecognitionAnswersTG(ctx)
			if err != nil {
				log.Println(err)
			}
			for _, result := range results {
				// Send answer image
				newFilename := fmt.Sprintf(result.ResImgPathName)
				photoBytes, err := os.ReadFile(newFilename)
				if err != nil {
					log.Println(err)
				}

				photoFileBytes := tgbotapi.FileBytes{
					Name:  "picture",
					Bytes: photoBytes,
				}

				_, err = tg.BotAPI.Send(tgbotapi.NewPhoto(result.ChatID, photoFileBytes))
				if err != nil {
					log.Println(err)
				}
				// Send result text
				description := tg.describe(result.Description)
				msg := tgbotapi.NewMessage(result.ChatID, description)
				_, err = tg.BotAPI.Send(msg)
				if err != nil {
					log.Println(err)
				}

			}
		}
	}

}

// Answer - send images with detection and description to them
func (tg *Telegram) describe(descriptions string) string {
	re, _ := regexp.Compile(`class...(\d+)`)
	res := re.FindAllStringSubmatch(descriptions, -1)
	nums := make([]int, len(res))
	str := ""
	for i, match := range res {
		num, err := strconv.Atoi(match[1])
		if err != nil {
			panic(err)
		}
		nums[i] = num
		str = str + tg.classes.Classes[num]
	}

	return str
}

// Функция для загрузки файла с указанного URL
func (tg *Telegram) downloadFile(url string, messageID string) (string, error) {
	date := time.Now().Format(layout)
	fileName := date + "/to_detect" + messageID + ".jpg"
	filePathAndName := fmt.Sprintf(tg.filePath+"%s", fileName)
	// Создаем файл на сервере
	if _, err := os.Stat(tg.filePath + "/" + date); os.IsNotExist(err) {
		os.MkdirAll(tg.filePath+"/"+date, 0700) // Create your file
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
	if err != nil {
		log.Println(err)
	}
	// Resize to 640x640
	newImage := resize.Resize(640, 640, receivedImage, resize.Lanczos3)
	// Save image
	err = jpeg.Encode(out, newImage, nil)
	if err != nil {
		return "", err
	}
	return filePathAndName, nil
}
