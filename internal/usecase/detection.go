package usecase

import (
	"context"
	"github.com/evrone/go-clean-template/internal/entity"
)

type RecognitionUseCase struct {
	Repo Repo
	Auth Auth
}

func New(repo Repo, auth Auth) *RecognitionUseCase {
	return &RecognitionUseCase{
		Repo: repo,
		Auth: auth,
	}
}

func (ir *RecognitionUseCase) AddRequest(ctx context.Context, req entity.Request) error {
	if req.ChatID == 0 {
		return ir.Repo.AddRequestWEB(ctx, req)
	}
	return ir.Repo.AddRequest(ctx, req)
}

func (ir *RecognitionUseCase) GetRecognitionTasks(ctx context.Context) ([]entity.Request, error) {
	return ir.Repo.GetRecognitionTasks(ctx)
}

func (ir *RecognitionUseCase) MakeRecognized(ctx context.Context, req entity.Request) error {
	return ir.Repo.MakeRecognized(ctx, req)
}

func (ir *RecognitionUseCase) GetRecognitionAnswersTG(ctx context.Context) ([]entity.Request, error) {
	return ir.Repo.GetRecognitionAnswersTG(ctx)
}

func (ir *RecognitionUseCase) GetRecognitionAnswersWEB(ctx context.Context, userName string) ([]entity.Request, error) {
	return ir.Repo.GetRecognitionAnswersWEB(ctx, userName)
}

func (ir *RecognitionUseCase) AddUser(c context.Context, u entity.User) error {
	return ir.Auth.AddUser(c, u)
}

func (ir *RecognitionUseCase) Login(c context.Context, u entity.User) (bool, error) {
	return ir.Auth.Login(c, u)
}
