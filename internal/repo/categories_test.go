package repo

import (
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func (s *pgdbTestSuite) TestCategoriesRepo_Create() {
	news := s.createNews()

	testCases := []struct {
		testName   string
		newsId     int64
		categories []int64
		expectErr  error
	}{
		{
			testName:   "correct test",
			newsId:     news.Id,
			categories: []int64{1, 2, 3, 4},
			expectErr:  nil,
		},
		{
			testName:   "news does not exists",
			newsId:     -1,
			categories: []int64{1, 2, 3, 4},
			expectErr:  ErrNotFound,
		},
		{
			testName:   "categories already exists",
			newsId:     news.Id,
			categories: []int64{1, 2, 2},
			expectErr:  ErrAlreadyExists,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.testName, func(t *testing.T) {
			err := s.categories.Create(s.tx.DB(s.ctx), tc.newsId, tc.categories)

			assert.Equal(t, tc.expectErr, err)

			if tc.expectErr == nil {
				sql := "SELECT category_id FROM news_categories WHERE news_id = $1"

				rows, err := s.pg.Query(s.ctx, sql, tc.newsId)

				assert.NoError(t, err)

				actual, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (id int64, err error) {
					err = row.Scan(&id)
					return
				})

				assert.NoError(t, err)
				assert.Equal(t, tc.categories, actual)
			}
		})
	}
}

func (s *pgdbTestSuite) TestCategoriesRepo_Delete() {
	news := s.createNews()
	s.createCategories(news.Id, 123)

	testCases := []struct {
		testName  string
		newsId    int64
		expectErr error
	}{
		{
			testName:  "correct test",
			newsId:    news.Id,
			expectErr: nil,
		},
		{
			testName:  "categories not found",
			newsId:    -1,
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.testName, func(t *testing.T) {
			err := s.categories.Delete(s.tx.DB(s.ctx), tc.newsId)

			assert.Equal(t, tc.expectErr, err)

			if tc.expectErr == nil {
				sql := "SELECT EXISTS(SELECT 1 FROM news_categories WHERE news_id = $1)"

				var exists bool
				err = s.pg.QueryRow(s.ctx, sql, tc.newsId).Scan(&exists)
				assert.NoError(t, err)

				assert.False(t, exists)
			}
		})
	}
}
