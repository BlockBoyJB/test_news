package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"test_news/internal/mocks/repomocks"
	"test_news/internal/mocks/txmocks"
	"test_news/internal/model"
	"test_news/internal/repo"
	"test_news/internal/repo/txmanager"
	"testing"
)

var errUnexpectedError = errors.New("some error")

func TestNewsService_Create(t *testing.T) {
	type args struct {
		ctx   context.Context
		input model.News
	}

	type mockBehaviour func(
		n *repomocks.MockNews,
		c *repomocks.MockCategories,
		mgr *txmocks.MockManager,
		tx *txmocks.MockTX,
		a args,
	)

	testCases := []struct {
		testName      string
		args          args
		mockBehaviour mockBehaviour
		expectOutput  int64
		expectErr     error
	}{
		{
			testName: "correct test",
			args: args{
				ctx: context.Background(),
				input: model.News{
					Title:      "FOOBAR",
					Content:    "CONTENT",
					Categories: nil,
				},
			},
			mockBehaviour: func(n *repomocks.MockNews, c *repomocks.MockCategories, mgr *txmocks.MockManager, tx *txmocks.MockTX, a args) {
				mgr.EXPECT().TxFunc(a.ctx, gomock.Any()).DoAndReturn(mockTX(tx))

				n.EXPECT().Create(tx, model.News{
					Title:   a.input.Title,
					Content: a.input.Content,
				}).Return(int64(1), nil)
				tx.EXPECT().Commit(a.ctx).Return(nil)
			},
			expectOutput: 1,
			expectErr:    nil,
		},
		{
			testName: "correct test with categories",
			args: args{
				ctx: context.Background(),
				input: model.News{
					Title:      "FOOBAR",
					Content:    "CONTENT",
					Categories: []int64{1, 2, 3, 4},
				},
			},
			mockBehaviour: func(n *repomocks.MockNews, c *repomocks.MockCategories, mgr *txmocks.MockManager, tx *txmocks.MockTX, a args) {
				mgr.EXPECT().TxFunc(a.ctx, gomock.Any()).DoAndReturn(mockTX(tx))
				n.EXPECT().Create(tx, model.News{
					Title:      a.input.Title,
					Content:    a.input.Content,
					Categories: a.input.Categories,
				}).Return(int64(1), nil)
				c.EXPECT().Create(tx, int64(1), a.input.Categories).Return(nil)
				tx.EXPECT().Commit(a.ctx).Return(nil)
			},
			expectOutput: 1,
			expectErr:    nil,
		},
		{
			testName: "unexpected error create user",
			args: args{
				ctx: context.Background(),
				input: model.News{
					Title:      "FOOBAR",
					Content:    "CONTENT",
					Categories: []int64{1, 2, 3, 4},
				},
			},
			mockBehaviour: func(n *repomocks.MockNews, c *repomocks.MockCategories, mgr *txmocks.MockManager, tx *txmocks.MockTX, a args) {
				mgr.EXPECT().TxFunc(a.ctx, gomock.Any()).DoAndReturn(mockTX(tx))
				n.EXPECT().Create(tx, model.News{
					Title:      a.input.Title,
					Content:    a.input.Content,
					Categories: a.input.Categories,
				}).Return(int64(0), errUnexpectedError)
				tx.EXPECT().Rollback(a.ctx).Return(nil)
			},
			expectOutput: 0,
			expectErr:    errUnexpectedError,
		},
		{
			testName: "categories already exists",
			args: args{
				ctx: context.Background(),
				input: model.News{
					Title:      "FOOBAR",
					Content:    "CONTENT",
					Categories: []int64{1, 2, 2},
				},
			},
			mockBehaviour: func(n *repomocks.MockNews, c *repomocks.MockCategories, mgr *txmocks.MockManager, tx *txmocks.MockTX, a args) {
				mgr.EXPECT().TxFunc(a.ctx, gomock.Any()).DoAndReturn(mockTX(tx))
				n.EXPECT().Create(tx, model.News{
					Title:      a.input.Title,
					Content:    a.input.Content,
					Categories: a.input.Categories,
				}).Return(int64(1), nil)
				c.EXPECT().Create(tx, int64(1), a.input.Categories).Return(repo.ErrAlreadyExists)
				tx.EXPECT().Rollback(a.ctx).Return(nil)
			},
			expectOutput: 0,
			expectErr:    ErrCategoriesAlreadyExists,
		},
		{
			testName: "commit error",
			args: args{
				ctx: context.Background(),
				input: model.News{
					Title:      "FOOBAR",
					Content:    "CONTENT",
					Categories: []int64{1, 2, 3, 4},
				},
			},
			mockBehaviour: func(n *repomocks.MockNews, c *repomocks.MockCategories, mgr *txmocks.MockManager, tx *txmocks.MockTX, a args) {
				mgr.EXPECT().TxFunc(a.ctx, gomock.Any()).DoAndReturn(mockTX(tx))
				n.EXPECT().Create(tx, model.News{
					Title:      a.input.Title,
					Content:    a.input.Content,
					Categories: a.input.Categories,
				}).Return(int64(1), nil)
				c.EXPECT().Create(tx, int64(1), a.input.Categories).Return(nil)
				tx.EXPECT().Commit(a.ctx).Return(errUnexpectedError)
				tx.EXPECT().Rollback(a.ctx).Return(nil)
			},
			expectOutput: 0,
			expectErr:    errUnexpectedError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			n := repomocks.NewMockNews(ctrl)
			c := repomocks.NewMockCategories(ctrl)
			mgr := txmocks.NewMockManager(ctrl)
			tx := txmocks.NewMockTX(ctrl)

			tc.mockBehaviour(n, c, mgr, tx, tc.args)

			s := newNewsService(mgr, n, c)

			id, err := s.Create(tc.args.ctx, tc.args.input)

			assert.ErrorIs(t, err, tc.expectErr)
			assert.Equal(t, tc.expectOutput, id)
		})
	}
}

func TestNewsService_Update(t *testing.T) {
	type args struct {
		ctx   context.Context
		input NewsUpdate
	}

	type mockBehaviour func(
		n *repomocks.MockNews,
		c *repomocks.MockCategories,
		mgr *txmocks.MockManager,
		tx *txmocks.MockTX,
		a args,
	)

	testCases := []struct {
		testName      string
		args          args
		mockBehaviour mockBehaviour
		expectErr     error
	}{
		{
			testName: "correct test",
			args: args{
				ctx: context.Background(),
				input: NewsUpdate{
					Id:         1,
					Title:      ptr("Foobar"),
					Content:    ptr("New content"),
					Categories: nil,
				},
			},
			mockBehaviour: func(n *repomocks.MockNews, c *repomocks.MockCategories, mgr *txmocks.MockManager, tx *txmocks.MockTX, a args) {
				mgr.EXPECT().TxFunc(a.ctx, gomock.Any()).DoAndReturn(mockTX(tx))

				n.EXPECT().Update(tx, a.input.Id, a.input.Title, a.input.Content).Return(nil)
				tx.EXPECT().Commit(a.ctx).Return(nil)
			},
			expectErr: nil,
		},
		{
			testName: "correct test with categories",
			args: args{
				ctx: context.Background(),
				input: NewsUpdate{
					Id:         1,
					Title:      ptr("Foobar"),
					Content:    ptr("New content"),
					Categories: []int64{1, 2, 3},
				},
			},
			mockBehaviour: func(n *repomocks.MockNews, c *repomocks.MockCategories, mgr *txmocks.MockManager, tx *txmocks.MockTX, a args) {
				mgr.EXPECT().TxFunc(a.ctx, gomock.Any()).DoAndReturn(mockTX(tx))

				n.EXPECT().Update(tx, a.input.Id, a.input.Title, a.input.Content).Return(nil)
				c.EXPECT().Delete(tx, a.input.Id).Return(nil)
				c.EXPECT().Create(tx, a.input.Id, a.input.Categories).Return(nil)
				tx.EXPECT().Commit(a.ctx).Return(nil)
			},
			expectErr: nil,
		},
		{
			testName: "news not found",
			args: args{
				ctx: context.Background(),
				input: NewsUpdate{
					Id:         0,
					Title:      ptr("Foobar"),
					Content:    ptr("New content"),
					Categories: nil,
				},
			},
			mockBehaviour: func(n *repomocks.MockNews, c *repomocks.MockCategories, mgr *txmocks.MockManager, tx *txmocks.MockTX, a args) {
				mgr.EXPECT().TxFunc(a.ctx, gomock.Any()).DoAndReturn(mockTX(tx))

				n.EXPECT().Update(tx, a.input.Id, a.input.Title, a.input.Content).Return(repo.ErrNotFound)
				tx.EXPECT().Rollback(a.ctx).Return(nil)
			},
			expectErr: ErrNewsNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			n := repomocks.NewMockNews(ctrl)
			c := repomocks.NewMockCategories(ctrl)
			mgr := txmocks.NewMockManager(ctrl)
			tx := txmocks.NewMockTX(ctrl)

			tc.mockBehaviour(n, c, mgr, tx, tc.args)

			s := newNewsService(mgr, n, c)

			err := s.Update(tc.args.ctx, tc.args.input)

			if !errors.Is(err, tc.expectErr) {
				t.Fail()
			}
		})
	}
}

func TestNewsService_FindWithCategories(t *testing.T) {
	type args struct {
		ctx    context.Context
		limit  int
		offset int
	}

	type mockBehaviour func(
		n *repomocks.MockNews,
		mgr *txmocks.MockManager,
		exec *txmocks.MockExecutor,
		a args,
	)

	testCases := []struct {
		testName      string
		args          args
		mockBehaviour mockBehaviour
		expectOutput  []model.News
		expectErr     error
	}{
		{
			testName: "correct test",
			args: args{
				ctx:    context.Background(),
				limit:  20,
				offset: 0,
			},
			mockBehaviour: func(n *repomocks.MockNews, mgr *txmocks.MockManager, exec *txmocks.MockExecutor, a args) {
				mgr.EXPECT().DB(a.ctx).Return(exec)
				n.EXPECT().FindWithCategories(exec, a.limit, a.offset).Return([]model.News{
					{
						Id:         1,
						Title:      "Title 1",
						Content:    "Content 1",
						Categories: []int64{1, 2, 3},
					},
					{
						Id:         2,
						Title:      "Title 2",
						Content:    "Content 2",
						Categories: []int64{2, 4},
					},
				}, nil)
			},
			expectOutput: []model.News{
				{
					Id:         1,
					Title:      "Title 1",
					Content:    "Content 1",
					Categories: []int64{1, 2, 3},
				},
				{
					Id:         2,
					Title:      "Title 2",
					Content:    "Content 2",
					Categories: []int64{2, 4},
				},
			},
			expectErr: nil,
		},
		{
			testName: "unexpected error",
			args: args{
				ctx:    context.Background(),
				limit:  20,
				offset: 0,
			},
			mockBehaviour: func(n *repomocks.MockNews, mgr *txmocks.MockManager, exec *txmocks.MockExecutor, a args) {
				mgr.EXPECT().DB(a.ctx).Return(exec)
				n.EXPECT().FindWithCategories(exec, a.limit, a.offset).Return(nil, errUnexpectedError)
			},
			expectOutput: nil,
			expectErr:    errUnexpectedError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			n := repomocks.NewMockNews(ctrl)
			mgr := txmocks.NewMockManager(ctrl)
			exec := txmocks.NewMockExecutor(ctrl)

			tc.mockBehaviour(n, mgr, exec, tc.args)

			s := newNewsService(mgr, n, nil)

			actual, err := s.FindWithCategories(tc.args.ctx, tc.args.limit, tc.args.offset)

			if !errors.Is(err, tc.expectErr) {
				t.Fail()
			}

			assert.Equal(t, tc.expectOutput, actual)
		})
	}
}

func mockTX(tx *txmocks.MockTX) func(ctx context.Context, f func(tx txmanager.TX) error) (err error) {
	return func(ctx context.Context, f func(tx txmanager.TX) error) (err error) {
		if f == nil {
			return errors.New("nil tx func")
		}
		defer func() {
			if tx == nil {
				return
			}
			if e := tx.Rollback(ctx); e != nil {
				err = e
			}
		}()
		if err = f(tx); err != nil {
			return
		}
		if err = tx.Commit(ctx); err != nil {
			return
		}
		tx = nil
		return nil
	}
}

func ptr[T any](t T) *T {
	return &t
}
