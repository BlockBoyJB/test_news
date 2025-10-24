package service

import (
	"context"
	"errors"
	"fmt"
	"test_news/internal/model"
	"test_news/internal/repo"
	"test_news/internal/repo/txmanager"
)

type newsService struct {
	tx         txmanager.Manager
	news       repo.News
	categories repo.Categories
}

func newNewsService(tx txmanager.Manager, news repo.News, categories repo.Categories) *newsService {
	return &newsService{
		tx:         tx,
		news:       news,
		categories: categories,
	}
}

func (s *newsService) Create(ctx context.Context, news model.News) (int64, error) {
	const op = "service.news.Create"

	var result int64
	err := s.tx.TxFunc(ctx, func(tx txmanager.TX) error {
		id, err := s.news.Create(tx, news)
		if err != nil {
			return fmt.Errorf("%s create user error: %w", op, err)
		}
		// на всякий случай проверка
		if len(news.Categories) != 0 {
			if err = s.categories.Create(tx, id, news.Categories); err != nil {
				if errors.Is(err, repo.ErrAlreadyExists) {
					return ErrCategoriesAlreadyExists
				}
				return fmt.Errorf("%s create categories error: %w", op, err)
			}
		}
		result = id
		return nil
	})
	if err != nil {
		return 0, err
	}
	return result, nil
}

type NewsUpdate struct {
	Id         int64
	Title      *string
	Content    *string
	Categories []int64
}

func (s *newsService) Update(ctx context.Context, input NewsUpdate) error {
	const op = "service.news.Update"

	err := s.tx.TxFunc(ctx, func(tx txmanager.TX) error {
		err := s.news.Update(tx, input.Id, input.Title, input.Content)
		if err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				return ErrNewsNotFound
			}
			return fmt.Errorf("%s update news error: %w", op, err)
		}

		if len(input.Categories) != 0 {
			if err = s.categories.Delete(tx, input.Id); err != nil {
				return fmt.Errorf("%s delete categories error: %w", op, err)
			}

			if err = s.categories.Create(tx, input.Id, input.Categories); err != nil {
				if errors.Is(err, repo.ErrAlreadyExists) {
					return ErrCategoriesAlreadyExists
				}
				return fmt.Errorf("%s create categories error: %w", op, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *newsService) FindWithCategories(ctx context.Context, limit, offset int) ([]model.News, error) {
	const op = "service.news.FindWithCategories"

	news, err := s.news.FindWithCategories(s.tx.DB(ctx), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return news, nil
}
