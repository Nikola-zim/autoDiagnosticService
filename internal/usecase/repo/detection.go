package repo

import (
	"context"
	"fmt"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/jackc/pgx/v4"
	"strings"
)

// RecognitionRepo -.
type RecognitionRepo struct {
	postgres  *postgres.Postgres
	batchSize int16
}

const defaultBatchSize = 32

// NewRecognitionRepo -.
func NewRecognitionRepo(pg *postgres.Postgres) *RecognitionRepo {
	return &RecognitionRepo{
		postgres:  pg,
		batchSize: defaultBatchSize,
	}
}

func (ir *RecognitionRepo) AddRequest(ctx context.Context, request entity.Request) error {
	// Start transaction
	tx, err := ir.postgres.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
	}(tx, ctx)
	// Insert request
	_, err = tx.Exec(ctx, newTask,
		request.ChatID, request.ImagePathName, Ready)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (ir *RecognitionRepo) GetRecognitionTasks(ctx context.Context) ([]entity.Request, error) {
	// Start transaction
	tx, err := ir.postgres.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	// Defer a rollback in case anything fails.
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
	}(tx, ctx)
	// Get tasks
	rows, err := tx.Query(ctx, getTasks, Ready, ir.batchSize)

	tasks := make([]entity.Request, 0, ir.batchSize)
	for rows.Next() {
		var task entity.Request
		err = rows.Scan(&task.ID, &task.ChatID, &task.ImagePathName)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (ir *RecognitionRepo) MakeRecognized(ctx context.Context, req entity.Request) error {
	// Start transaction
	tx, err := ir.postgres.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
	}(tx, ctx)

	//
	_, err = tx.Exec(ctx, changeToRecognized, req.ResImgPathName, req.Description, Recognized, req.ID)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (ir *RecognitionRepo) GetRecognitionAnswers(ctx context.Context) ([]entity.Request, error) {
	// Start transaction
	tx, err := ir.postgres.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	// Defer a rollback in case anything fails.
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
	}(tx, ctx)

	results := make([]entity.Request, 0, ir.batchSize)
	allID := make([]interface{}, 0, ir.batchSize)
	// Get answers
	rows, err := tx.Query(ctx, getAnswers)
	valueStrings := make([]string, 0, ir.batchSize)
	for rows.Next() {
		iter := 1
		var result entity.Request
		err = rows.Scan(&result.ID, &result.ChatID, &result.ResImgPathName, &result.Description)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
		valueStrings = append(valueStrings, fmt.Sprintf("$%d", iter))
		iter++
		allID = append(allID, int32(result.ID))
	}
	if len(allID) > 0 {
		query := fmt.Sprintf(changeToDone, strings.Join(valueStrings, ","))
		_, err = tx.Exec(ctx, query, allID...)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	return results, nil
}
