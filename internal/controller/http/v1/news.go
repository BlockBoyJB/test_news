package v1

import (
	"github.com/gofiber/fiber/v3"
	"strconv"
	"test_news/internal/model"
	"test_news/internal/service"
)

type newsRouter struct {
	news service.News
}

func newNewsRouter(g fiber.Router, news service.News) {
	r := &newsRouter{
		news: news,
	}

	g.Post("/create", r.create)
	g.Post("/edit/:id", r.update)
	g.Get("/list", r.list)
}

type newsCreateInput struct {
	Title      string  `json:"Title" validate:"required"`
	Content    string  `json:"Content" validate:"required"`
	Categories []int64 `json:"Categories" validate:"required,dive,gt=0"`
}

func (r *newsRouter) create(c fiber.Ctx) error {
	var input newsCreateInput

	if err := c.Bind().Body(&input); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	id, err := r.news.Create(c.Context(), model.News{
		Title:      input.Title,
		Content:    input.Content,
		Categories: input.Categories,
	})
	if err != nil {
		return err
	}
	response := struct {
		Id int64 `json:"Id"`
	}{
		Id: id,
	}
	return c.JSON(response)
}

type newsUpdateInput struct {
	Title      *string `json:"Title"`
	Content    *string `json:"Content"`
	Categories []int64 `json:"Categories" validate:"dive,gt=0"`
}

func (r *newsRouter) update(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var input newsUpdateInput

	if err = c.Bind().Body(&input); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = r.news.Update(c.Context(), service.NewsUpdate{
		Id:         int64(id),
		Title:      input.Title,
		Content:    input.Content,
		Categories: input.Categories,
	})
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

type newsPaginationInput struct {
	Limit  int `query:"limit" validate:"gte=0,lte=20"`
	Offset int `query:"offset" validate:"gte=0"`
}

type newsListResponse struct {
	Success bool         `json:"Success"`
	News    []model.News `json:"News"`
}

func (r *newsRouter) list(c fiber.Ctx) error {
	var input newsPaginationInput

	if err := c.Bind().Query(&input); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	news, err := r.news.FindWithCategories(c.Context(), input.Limit, input.Offset)
	if err != nil {
		return err
	}

	return c.JSON(newsListResponse{
		Success: true,
		News:    news,
	})
}
