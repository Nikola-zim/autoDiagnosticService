package usecase

import (
	"context"
	"github.com/evrone/go-clean-template/internal/entity"
)

type ImageRecognitionUseCase struct {
	Repo ImagesRepo
}

func New(repo ImagesRepo) *ImageRecognitionUseCase {
	return &ImageRecognitionUseCase{
		Repo: repo,
	}
}

func (ir *ImageRecognitionUseCase) AddRequest(ctx context.Context, req entity.Request) error {
	return ir.Repo.AddRequest(ctx, req)
}

func (ir *ImageRecognitionUseCase) GetRecognitionTasks(ctx context.Context) ([]entity.Request, error) {
	return ir.Repo.GetRecognitionTasks(ctx)
}

func (ir *ImageRecognitionUseCase) MakeRecognized(ctx context.Context, req entity.Request) error {
	return ir.Repo.MakeRecognized(ctx, req)
}

func (ir *ImageRecognitionUseCase) GetRecognitionAnswers(ctx context.Context) ([]entity.Request, error) {
	return ir.Repo.GetRecognitionAnswers(ctx)
}
