package v1

import (
	"github.com/gofiber/fiber/v3"
	"test_news/internal/service"
)

func NewRouter(g fiber.Router, services *service.Services) {
	g.Get("/ping", ping)
	g.Get("/authorize", createTokenHandler(services.Auth))

	g.Use(errorMiddleware)

	v1 := g.Group("/api/v1", authMiddleware(services.Auth))
	newNewsRouter(v1.Group("/news"), services.News)
}

func ping(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func createTokenHandler(auth service.Auth) fiber.Handler {
	return func(c fiber.Ctx) error {
		token, err := auth.Create()
		if err != nil {
			return err
		}
		response := struct {
			Token string `json:"token"`
		}{
			Token: token,
		}
		return c.JSON(response)
	}
}
