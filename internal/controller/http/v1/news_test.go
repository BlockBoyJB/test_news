package v1

import (
	"bytes"
	"context"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"test_news/internal/mocks/servicemocks"
	"test_news/internal/model"
	"test_news/internal/service"
	"test_news/pkg/validator"
	"testing"
)

func TestNewsRouter_create(t *testing.T) {
	type args struct {
		ctx   context.Context
		input model.News
	}

	type mockBehaviour func(n *servicemocks.MockNews, a args)

	testCases := []struct {
		testName      string
		args          args
		mockBehaviour mockBehaviour
		inputBody     string
		expectCode    int
		expectBody    string
	}{
		{
			testName: "correct test",
			args: args{
				ctx: context.Background(),
				input: model.News{
					Title:      "hello world",
					Content:    "my content",
					Categories: []int64{1, 2, 3},
				},
			},
			mockBehaviour: func(n *servicemocks.MockNews, a args) {
				n.EXPECT().Create(a.ctx, a.input).Return(int64(1), nil)
			},
			inputBody:  `{"Title": "hello world", "Content": "my content", "Categories": [1, 2, 3]}`,
			expectCode: fiber.StatusOK,
			expectBody: `{"Id":1}`,
		},
		{
			testName:      "missing field title",
			mockBehaviour: func(n *servicemocks.MockNews, a args) {},
			inputBody:     `{"Content": "my content", "Categories": [1, 2, 3]}`,
			expectCode:    fiber.StatusBadRequest,
			expectBody:    fiber.ErrBadRequest.Message,
		},
		{
			testName:      "missing field content",
			mockBehaviour: func(n *servicemocks.MockNews, a args) {},
			inputBody:     `{"Title": "hello world", "Categories": [1, 2, 3]}`,
			expectCode:    fiber.StatusBadRequest,
			expectBody:    fiber.ErrBadRequest.Message,
		},
		{
			testName:      "missing field categories",
			mockBehaviour: func(n *servicemocks.MockNews, a args) {},
			inputBody:     `{"Title": "hello world", "Content": "my content"}`,
			expectCode:    fiber.StatusBadRequest,
			expectBody:    fiber.ErrBadRequest.Message,
		},
		{
			testName:      "categories values <= 0",
			mockBehaviour: func(n *servicemocks.MockNews, a args) {},
			inputBody:     `{"Title": "hello world", "Content": "my content", "Categories": [0, -1, 3]}`,
			expectCode:    fiber.StatusBadRequest,
			expectBody:    fiber.ErrBadRequest.Message,
		},
		{
			testName: "categories already exists",
			args: args{
				ctx: context.Background(),
				input: model.News{
					Title:      "hello world",
					Content:    "my content",
					Categories: []int64{1, 1},
				},
			},
			mockBehaviour: func(n *servicemocks.MockNews, a args) {
				n.EXPECT().Create(a.ctx, a.input).Return(int64(0), service.ErrCategoriesAlreadyExists)
			},
			inputBody:  `{"Title": "hello world", "Content": "my content", "Categories": [1, 1]}`,
			expectCode: fiber.StatusBadRequest,
			expectBody: service.ErrCategoriesAlreadyExists.Error(),
		},
		{
			testName: "unexpected error",
			args: args{
				ctx: context.Background(),
				input: model.News{
					Title:      "hello world",
					Content:    "my content",
					Categories: []int64{1, 2, 3},
				},
			},
			mockBehaviour: func(n *servicemocks.MockNews, a args) {
				n.EXPECT().Create(a.ctx, a.input).Return(int64(0), errors.New("some error"))
			},
			inputBody:  `{"Title": "hello world", "Content": "my content", "Categories": [1, 2, 3]}`,
			expectCode: fiber.StatusInternalServerError,
			expectBody: fiber.ErrInternalServerError.Message,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			n := servicemocks.NewMockNews(ctrl)
			a := servicemocks.NewMockAuth(ctrl)

			a.EXPECT().Validate("TOKEN").Return(true)

			tc.mockBehaviour(n, tc.args)

			h := fiber.New(fiber.Config{
				StructValidator: validator.New(),
			})
			NewRouter(h, &service.Services{
				Auth: a,
				News: n,
			})

			r := httptest.NewRequest(fiber.MethodPost, "/api/v1/news/create", bytes.NewBufferString(tc.inputBody))

			r.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			r.Header.Set(fiber.HeaderAuthorization, "Bearer TOKEN")

			resp, err := h.Test(r)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectBody, string(body))
		})
	}
}

func TestNewsRouter_update(t *testing.T) {
	type args struct {
		ctx   context.Context
		input service.NewsUpdate
	}

	type mockBehaviour func(n *servicemocks.MockNews, a args)

	testCases := []struct {
		testName      string
		args          args
		mockBehaviour mockBehaviour
		inputBody     string
		inputId       string
		expectCode    int
		expectBody    string
	}{
		{
			testName: "correct test",
			args: args{
				ctx: context.Background(),
				input: service.NewsUpdate{
					Id:         1,
					Title:      ptr("new world"),
					Content:    ptr("new content"),
					Categories: []int64{1, 2, 3},
				},
			},
			mockBehaviour: func(n *servicemocks.MockNews, a args) {
				n.EXPECT().Update(a.ctx, a.input).Return(nil)
			},
			inputBody:  `{"Title": "new world", "Content": "new content", "Categories": [1, 2, 3]}`,
			inputId:    "1",
			expectCode: fiber.StatusOK,
			expectBody: "OK",
		},
		{
			testName:      "incorrect id",
			mockBehaviour: func(n *servicemocks.MockNews, a args) {},
			inputBody:     `{"Content": "my content", "Categories": [1, 2, 3]}`,
			inputId:       "foobar",
			expectCode:    fiber.StatusBadRequest,
			expectBody:    fiber.ErrBadRequest.Message,
		},
		{
			testName:      "categories values <= 0",
			mockBehaviour: func(n *servicemocks.MockNews, a args) {},
			inputBody:     `{"Title": "hello world", "Content": "my content", "Categories": [0, -1, 3]}`,
			inputId:       "1",
			expectCode:    fiber.StatusBadRequest,
			expectBody:    fiber.ErrBadRequest.Message,
		},
		{
			testName: "news not found",
			args: args{
				ctx: context.Background(),
				input: service.NewsUpdate{
					Id:         2,
					Title:      ptr("New title"),
					Content:    nil,
					Categories: nil,
				},
			},
			mockBehaviour: func(n *servicemocks.MockNews, a args) {
				n.EXPECT().Update(a.ctx, a.input).Return(service.ErrNewsNotFound)
			},
			inputBody:  `{"Title": "New title"}`,
			inputId:    "2",
			expectCode: fiber.StatusNotFound,
			expectBody: service.ErrNewsNotFound.Error(),
		},
		{
			testName: "categories already exists",
			args: args{
				ctx: context.Background(),
				input: service.NewsUpdate{
					Id:         1,
					Title:      nil,
					Content:    nil,
					Categories: []int64{1, 1},
				},
			},
			mockBehaviour: func(n *servicemocks.MockNews, a args) {
				n.EXPECT().Update(a.ctx, a.input).Return(service.ErrCategoriesAlreadyExists)
			},
			inputBody:  `{"Categories": [1, 1]}`,
			inputId:    "1",
			expectCode: fiber.StatusBadRequest,
			expectBody: service.ErrCategoriesAlreadyExists.Error(),
		},
		{
			testName: "unexpected error",
			args: args{
				ctx: context.Background(),
				input: service.NewsUpdate{
					Id:         1,
					Title:      ptr("hello world"),
					Categories: []int64{1, 2, 3},
				},
			},
			mockBehaviour: func(n *servicemocks.MockNews, a args) {
				n.EXPECT().Update(a.ctx, a.input).Return(errors.New("some error"))
			},
			inputBody:  `{"Title": "hello world", "Categories": [1, 2, 3]}`,
			inputId:    "1",
			expectCode: fiber.StatusInternalServerError,
			expectBody: fiber.ErrInternalServerError.Message,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			n := servicemocks.NewMockNews(ctrl)
			a := servicemocks.NewMockAuth(ctrl)

			a.EXPECT().Validate("TOKEN").Return(true)

			tc.mockBehaviour(n, tc.args)

			h := fiber.New(fiber.Config{
				StructValidator: validator.New(),
			})
			NewRouter(h, &service.Services{
				Auth: a,
				News: n,
			})

			r := httptest.NewRequest(fiber.MethodPost, "/api/v1/news/edit/"+tc.inputId, bytes.NewBufferString(tc.inputBody))

			r.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			r.Header.Set(fiber.HeaderAuthorization, "Bearer TOKEN")

			resp, err := h.Test(r)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectBody, string(body))
		})
	}
}

func TestNewsRouter_find(t *testing.T) {
	type args struct {
		ctx    context.Context
		limit  int
		offset int
	}

	type mockBehaviour func(n *servicemocks.MockNews, a args)

	testCases := []struct {
		testName      string
		args          args
		mockBehaviour mockBehaviour
		inputQuery    string
		expectBody    string
	}{
		{
			testName: "correct test",
			args: args{
				ctx:    context.Background(),
				limit:  10,
				offset: 0,
			},
			mockBehaviour: func(n *servicemocks.MockNews, a args) {
				n.EXPECT().FindWithCategories(a.ctx, a.limit, a.offset).Return([]model.News{
					{
						Id:         1,
						Title:      "Foobar",
						Content:    "Content",
						Categories: []int64{1, 2, 3},
					},
					{
						Id:         2,
						Title:      "Hello world",
						Content:    "Content",
						Categories: []int64{1},
					},
				}, nil)
			},
			inputQuery: `limit=10&offset=0`,
			expectBody: `{"Success":true,"News":[{"Id":1,"Title":"Foobar","Content":"Content","Categories":[1,2,3]},{"Id":2,"Title":"Hello world","Content":"Content","Categories":[1]}]}`,
		},
		{
			testName:      "incorrect limit #1",
			args:          args{},
			mockBehaviour: func(n *servicemocks.MockNews, a args) {},
			inputQuery:    `limit=-1&offset=0`,
			expectBody:    fiber.ErrBadRequest.Message,
		},
		{
			testName:      "incorrect limit #2",
			args:          args{},
			mockBehaviour: func(n *servicemocks.MockNews, a args) {},
			inputQuery:    `limit=23&offset=0`,
			expectBody:    fiber.ErrBadRequest.Message,
		},
		{
			testName:      "incorrect offset",
			args:          args{},
			mockBehaviour: func(n *servicemocks.MockNews, a args) {},
			inputQuery:    `limit=10&offset=-1`,
			expectBody:    fiber.ErrBadRequest.Message,
		},
		{
			testName: "unexpected error",
			args: args{
				ctx:    context.Background(),
				limit:  10,
				offset: 0,
			},
			mockBehaviour: func(n *servicemocks.MockNews, a args) {
				n.EXPECT().FindWithCategories(a.ctx, a.limit, a.offset).Return(nil, errors.New("some error"))
			},
			inputQuery: `limit=10&offset=0`,
			expectBody: fiber.ErrInternalServerError.Message,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			n := servicemocks.NewMockNews(ctrl)
			a := servicemocks.NewMockAuth(ctrl)

			a.EXPECT().Validate("TOKEN").Return(true)

			tc.mockBehaviour(n, tc.args)

			h := fiber.New(fiber.Config{
				StructValidator: validator.New(),
			})
			NewRouter(h, &service.Services{
				Auth: a,
				News: n,
			})

			r := httptest.NewRequest(fiber.MethodGet, "/api/v1/news/list?"+tc.inputQuery, nil)

			r.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			r.Header.Set(fiber.HeaderAuthorization, "Bearer TOKEN")

			resp, err := h.Test(r)
			assert.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectBody, string(body))
		})
	}
}

func ptr[T any](t T) *T {
	return &t
}
