package repo

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

type Categories interface {
	Create(exec Querier, newsId int64, categories []int64) error
	Delete(exec Querier, newsId int64) error
}

type categoriesRepo struct{}

func NewCategoriesRepo() Categories {
	return &categoriesRepo{}
}

func (r *categoriesRepo) Create(exec Querier, newsId int64, categories []int64) error {
	sql := "INSERT INTO news_categories (news_id, category_id) SELECT $1, unnest($2::bigint[])"
	if _, err := exec.Exec(sql, newsId, categories); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == codeErrForeignKeyViolation {
				return ErrNotFound
			}
			if pgErr.Code == codeErrUniqueViolation {
				return ErrAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (r *categoriesRepo) Delete(exec Querier, newsId int64) error {
	sql := "DELETE FROM news_categories WHERE news_id = $1"
	if _, err := exec.Exec(sql, newsId); err != nil {
		return err
	}
	return nil
}
