package v1

import (
	"autoDiagnosticService/internal/entity"
	"context"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/auth.go -package=mocks

type (
	// Recognition -.
	Recognition interface {
		AddRequest(context.Context, entity.Request) error
		GetRecognitionTasks(context.Context) ([]entity.Request, error)
		MakeRecognized(context.Context, entity.Request) error
		GetRecognitionAnswersTG(ctx context.Context) ([]entity.Request, error)
		GetRecognitionAnswersWEB(ctx context.Context, userName string) ([]entity.Request, error)
		Auth
		Balance
	}

	// Repo -.
	Repo interface {
		AddRequest(context.Context, entity.Request) error
		AddRequestWEB(context.Context, entity.Request) error
		GetRecognitionTasks(context.Context) ([]entity.Request, error)
		MakeRecognized(context.Context, entity.Request) error
		GetRecognitionAnswersTG(context.Context) ([]entity.Request, error)
		GetRecognitionAnswersWEB(ctx context.Context, userName string) ([]entity.Request, error)
		Balance
	}

	Balance interface {
		AddPoints(ctx context.Context, number int, userName string) error
	}

	Auth interface {
		AddUser(context.Context, entity.User) error
		Login(context.Context, entity.User) (bool, error)
	}
)
