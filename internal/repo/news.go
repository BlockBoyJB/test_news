package repo

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"strconv"
	"strings"
	"test_news/internal/model"
)

type News interface {
	Create(exec Querier, news model.News) (int64, error)
	Update(exec Querier, id int64, title, content *string) error
	FindWithCategories(exec Querier, limit, offset int) ([]model.News, error)
}

type newsRepo struct{}

func NewNewsRepo() News {
	return &newsRepo{}
}

func (r *newsRepo) Create(exec Querier, news model.News) (int64, error) {
	sql := "INSERT INTO news (title, content) VALUES ($1, $2) RETURNING id"

	var id int64
	if err := exec.QueryRow(sql, news.Title, news.Content).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *newsRepo) Update(exec Querier, id int64, title, content *string) error {
	var (
		parts []string
		args  []any
		pos   = 1
	)
	if title != nil {
		parts = append(parts, "title = $"+strconv.Itoa(pos))
		args = append(args, *title)
		pos++
	}
	if content != nil {
		parts = append(parts, "content = $"+strconv.Itoa(pos))
		args = append(args, *content)
		pos++
	}
	if len(parts) == 0 {
		return nil
	}
	sql := fmt.Sprintf("UPDATE news SET %s WHERE id = $%d", strings.Join(parts, ", "), pos)

	args = append(args, id)

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *newsRepo) FindWithCategories(exec Querier, limit, offset int) ([]model.News, error) {
	sql := `
		SELECT n.id, n.title, n.content, 
		       COALESCE(array_agg(nc.category_id) FILTER (WHERE nc.category_id IS NOT NULL), '{}') AS categories
		FROM news n
		LEFT JOIN news_categories nc ON n.id = nc.news_id
		GROUP BY n.id, n.title, n.content
		ORDER BY n.id
		LIMIT $1
		OFFSET $2
	`

	rows, err := exec.Query(sql, limit, offset)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[model.News])
}
