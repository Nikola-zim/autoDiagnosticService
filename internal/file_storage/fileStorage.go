package fileStorage

import (
	b64 "encoding/base64"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

const layout = "2006_01_02"

type FileStorage struct {
	StoragePath string
}

func New(storagePath string) *FileStorage {
	return &FileStorage{
		StoragePath: storagePath,
	}
}

func (fs *FileStorage) ImageSave(img string, pathName string) (string, error) {
	date := time.Now().Format(layout)
	decodedBytes, err := b64.StdEncoding.DecodeString(img)

	if err != nil {
		return "", err
	}

	// Создаем файл на сервере
	currentPath := fs.StoragePath + "/" + date
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
