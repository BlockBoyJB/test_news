package service

import (
	"context"
	"test_news/internal/model"
	"test_news/internal/repo"
	"test_news/internal/repo/txmanager"
)

type News interface {
	Create(ctx context.Context, news model.News) (int64, error)
	Update(ctx context.Context, input NewsUpdate) error
	FindWithCategories(ctx context.Context, limit, offset int) ([]model.News, error)
}

type Auth interface {
	Create() (string, error)
	Validate(tokenString string) bool
}

type (
	Services struct {
		Auth Auth
		News News
	}
	ServicesDependencies struct {
		NewsRepo       repo.News
		CategoriesRepo repo.Categories
		TxManager      txmanager.Manager
		JWTKey         string
	}
)

func NewServices(d *ServicesDependencies) *Services {
	return &Services{
		Auth: newAuthService(d.JWTKey),
		News: newNewsService(d.TxManager, d.NewsRepo, d.CategoriesRepo),
	}
}
