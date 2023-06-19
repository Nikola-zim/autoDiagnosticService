// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"github.com/evrone/go-clean-template/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// ImageRecognition -.
	ImageRecognition interface {
		AddRequest(context.Context, entity.Request) error
		GetRecognitionTasks(context.Context) ([]entity.Request, error)
		MakeRecognized(context.Context, entity.Request) error
		GetRecognitionAnswers(ctx context.Context) ([]entity.Request, error)
	}

	// ImagesRepo -.
	ImagesRepo interface {
		AddRequest(context.Context, entity.Request) error
		GetRecognitionTasks(context.Context) ([]entity.Request, error)
		MakeRecognized(context.Context, entity.Request) error
		GetRecognitionAnswers(context.Context) ([]entity.Request, error)
	}

	DetectionWorker interface {
	}
)
