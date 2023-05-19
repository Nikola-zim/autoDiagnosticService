package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Telegram struct {
	BotAPI   *tgbotapi.BotAPI
	botToken string
}

func New(token string, opts ...Option) (*Telegram, error) {

	botAPI, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, err
	}

	tg := &Telegram{
		botToken: token,
		BotAPI:   botAPI,
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

	for update := range updates {
		if update.Message != nil { // If we got a message
			//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			// Если сообщение содержит изображение
			if update.Message.Photo != nil {
				// Получаем информацию об изображении
				fileID := (update.Message.Photo)[len(update.Message.Photo)-1].FileID
				file, err := tg.BotAPI.GetFile(tgbotapi.FileConfig{FileID: fileID})
				if err != nil {
					log.Println(err)
					continue
				}

				fileName := strconv.Itoa(1234567) + "_" + strconv.Itoa(update.Message.MessageID) + ".jpg"
				filePathAndName := fmt.Sprintf("internal/controller/telegram/%s", fileName)
				// Загружаем изображение на сервер
				err = downloadFile(file.Link(tg.BotAPI.Token), filePathAndName)
				if err != nil {
					log.Println(err)
					continue
				}

				// Отправляем сообщение об успешной загрузке изображения
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Изображение успешно загружено на сервер.")
				_, err = tg.BotAPI.Send(msg)
				if err != nil {
					log.Println(err)
				}
				newFilename := fmt.Sprintf(filePathAndName)
				photoBytes, err := os.ReadFile(newFilename)
				if err != nil {
					panic(err)
				}

				photoFileBytes := tgbotapi.FileBytes{
					Name:  "picture",
					Bytes: photoBytes,
				}

				_, err = tg.BotAPI.Send(tgbotapi.NewPhoto(update.Message.Chat.ID, photoFileBytes))

				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "HIHIHIHIHI")
				//msg.ReplyToMessageID = update.Message.MessageID
				//
				//tg.BotAPI.Send(msg)
			}
		}
	}
}

// Функция для загрузки файла с указанного URL
func downloadFile(url string, fileName string) error {
	// Создаем файл на сервере
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// Получаем содержимое файла по URL и записываем в созданный файл
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
