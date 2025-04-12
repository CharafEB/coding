package database

import (
	"context"
	"database/sql"
	"github/codingMaster/internal/app/model"
)

type Store struct {
	Admin interface {
		Signup(ctx context.Context, val *model.User) error
		Login(ctx context.Context, val *model.Login) (*model.User, error)
	}

	Role interface {
		Unpprove(ctx context.Context, id int) error
		Upprove(ctx context.Context, id int, status string) error
		GetProject(ctx context.Context, id int) (*model.Project, error)
		AddProject(ctx context.Context, project *model.Project) error
		DeleteProject(ctx context.Context, id int) error
		UpdateProject(ctx context.Context, project *model.Project) error
	}
}

func NewStore(db *sql.DB) Store {
	if db == nil {
		panic("nil pointer passed to NewStore")
	}
	articleRepo := &articlesC{DB: db}

	return Store{
		Admin: articleRepo,
		Role:  articleRepo,
	}
}
