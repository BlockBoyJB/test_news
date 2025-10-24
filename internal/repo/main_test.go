package repo

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/suite"
	"test_news/internal/model"
	"test_news/internal/repo/txmanager"
	"test_news/pkg/postgres"
	"testing"
)

type pgdbTestSuite struct {
	suite.Suite
	ctx        context.Context
	pg         postgres.Postgres
	m          *migrate.Migrate
	tx         txmanager.Manager
	news       *newsRepo
	categories *categoriesRepo
}

func (s *pgdbTestSuite) SetupTest() {
	testPGUrl := "postgres://postgres:1234567890@localhost:6000/postgres"
	m, err := migrate.New("file://../../migrations", testPGUrl+"?sslmode=disable")
	if err != nil {
		panic(err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
	s.m = m

	s.ctx = context.Background()

	pg, err := postgres.NewPG(testPGUrl)
	if err != nil {
		panic(err)
	}
	s.pg = pg
	s.tx = txmanager.NewManager(pg)
	s.news = &newsRepo{}
	s.categories = &categoriesRepo{}
}

func (s *pgdbTestSuite) TearDownTest() {
	_ = s.m.Drop()
	s.pg.Close()
}

func TestPGDB(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(pgdbTestSuite))
}

func (s *pgdbTestSuite) createNews() model.News {
	n := model.News{
		Id:      0,
		Title:   "My Title",
		Content: "Content",
	}
	sql := "INSERT INTO news (title, content) VALUES ($1, $2) RETURNING id"
	var id int64
	if err := s.pg.QueryRow(s.ctx, sql, n.Title, n.Content).Scan(&id); err != nil {
		panic(err)
	}
	n.Id = id
	return n
}

func (s *pgdbTestSuite) createCategories(newsId int64, categoryId ...int64) {
	sql := "INSERT INTO news_categories (news_id, category_id) VALUES ($1, unnest($2::bigint[]))"
	if _, err := s.pg.Exec(s.ctx, sql, newsId, categoryId); err != nil {
		panic(err)
	}
}

func ptr[T any](t T) *T {
	return &t
}
