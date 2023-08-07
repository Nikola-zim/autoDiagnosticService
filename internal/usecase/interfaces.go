// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"github.com/evrone/go-clean-template/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Recognition -.
	Recognition interface {
		AddRequest(context.Context, entity.Request) error
		GetRecognitionTasks(context.Context) ([]entity.Request, error)
		MakeRecognized(context.Context, entity.Request) error
		GetRecognitionAnswersTG(ctx context.Context) ([]entity.Request, error)
		GetRecognitionAnswersWEB(ctx context.Context, userName string) ([]entity.Request, error)
		Auth
	}

	// Repo -.
	Repo interface {
		AddRequest(context.Context, entity.Request) error
		AddRequestWEB(context.Context, entity.Request) error
		GetRecognitionTasks(context.Context) ([]entity.Request, error)
		MakeRecognized(context.Context, entity.Request) error
		GetRecognitionAnswersTG(context.Context) ([]entity.Request, error)
		GetRecognitionAnswersWEB(ctx context.Context, userName string) ([]entity.Request, error)
	}

	Auth interface {
		AddUser(context.Context, entity.User) error
		Login(context.Context, entity.User) (bool, error)
	}

	DetectionWorker interface {
	}
)
