package repo

import (
	"github.com/stretchr/testify/assert"
	"test_news/internal/model"
	"testing"
)

func (s *pgdbTestSuite) TestNewsRepo_Create() {
	testCases := []struct {
		testName string
		news     model.News
	}{
		{
			testName: "correct test 1",
			news: model.News{
				Title:   "FOOBAR",
				Content: "TEXT 123123123123",
			},
		},
		{
			testName: "correct test 2",
			news: model.News{
				Title:   "HELLO WORLD",
				Content: "ANOTHER TEXT",
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.testName, func(t *testing.T) {
			id, err := s.news.Create(s.tx.DB(s.ctx), tc.news)

			assert.NoError(t, err)
			assert.True(t, id > 0) // pk cannot be <= 0

			sql := "SELECT title, content FROM news WHERE id = $1"
			var actual model.News

			err = s.pg.QueryRow(s.ctx, sql, id).Scan(&actual.Title, &actual.Content)

			assert.NoError(t, err)

			assert.Equal(t, tc.news, actual)
		})
	}
}

func (s *pgdbTestSuite) TestNewsRepo_Update() {
	news := s.createNews()

	testCases := []struct {
		testName     string
		id           int64
		title        *string
		content      *string
		expectUpdate model.News
		expectErr    error
	}{
		{
			testName: "update only title",
			id:       news.Id,
			title:    ptr("NEW TITLE"),
			expectUpdate: model.News{
				Title:   "NEW TITLE",
				Content: news.Content,
			},
			expectErr: nil,
		},
		{
			testName: "update only content",
			id:       news.Id,
			content:  ptr("NEW TEXT"),
			expectUpdate: model.News{
				Title:   "NEW TITLE",
				Content: "NEW TEXT",
			},
			expectErr: nil,
		},
		{
			testName: "update all",
			id:       news.Id,
			title:    ptr("NEW TITLE 2"),
			content:  ptr("NEW TEXT 2"),
			expectUpdate: model.News{
				Title:   "NEW TITLE 2",
				Content: "NEW TEXT 2",
			},
			expectErr: nil,
		},
		{
			testName:     "not found",
			id:           -1,
			title:        ptr("NEW TITLE 3"),
			content:      ptr("NEW TEXT 3"),
			expectUpdate: model.News{},
			expectErr:    ErrNotFound,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.testName, func(t *testing.T) {
			err := s.news.Update(s.tx.DB(s.ctx), tc.id, tc.title, tc.content)

			assert.Equal(t, tc.expectErr, err)

			if tc.expectErr == nil {
				sql := "SELECT title, content FROM news WHERE id = $1"

				var actual model.News
				err = s.pg.QueryRow(s.ctx, sql, tc.id).Scan(&actual.Title, &actual.Content)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectUpdate, actual)
			}
		})
	}
}

func (s *pgdbTestSuite) TestNewsRepo_FindWithCategories() {
	news1 := s.createNews()
	news2 := s.createNews()
	s.createCategories(news1.Id, 1, 2, 3, 4)
	s.createCategories(news2.Id, 2, 5)

	testCases := []struct {
		testName     string
		limit        int
		offset       int
		expectOutput []model.News
	}{
		{
			testName: "correct test",
			limit:    20,
			offset:   0,
			expectOutput: []model.News{
				{
					Id:         news1.Id,
					Title:      news1.Title,
					Content:    news1.Content,
					Categories: []int64{1, 2, 3, 4},
				},
				{
					Id:         news2.Id,
					Title:      news2.Title,
					Content:    news2.Content,
					Categories: []int64{2, 5},
				},
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.testName, func(t *testing.T) {
			news, err := s.news.FindWithCategories(s.tx.DB(s.ctx), tc.limit, tc.offset)

			assert.NoError(t, err)

			assert.Equal(t, tc.expectOutput, news)
		})
	}

}
