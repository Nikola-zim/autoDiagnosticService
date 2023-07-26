package repo

import (
	"context"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type Auth struct {
	postgres *postgres.Postgres
}

// NewAuth -.
func NewAuth(pg *postgres.Postgres) *Auth {
	return &Auth{
		postgres: pg,
	}
}

func (a *Auth) AddUser(ctx context.Context, user entity.User) error {
	// Insert request
	row := a.postgres.Pool.QueryRow(ctx, addUser, user.Login, user.Password)
	var userID int64
	err := row.Scan(&userID)
	if err != nil {
		return err
	}
	return nil
}

func (a *Auth) Login(ctx context.Context, user entity.User) (bool, error) {
	// Insert request
	row := a.postgres.Pool.QueryRow(ctx, login, user.Login, user.Password)
	var (
		userID int64
		exist  int64
	)
	err := row.Scan(&userID, &exist)
	if err != nil || exist == 0 {
		return false, err
	}
	return true, nil
}
