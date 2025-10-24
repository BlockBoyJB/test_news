package v1

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"strings"
	"test_news/internal/service"
)

func errorMiddleware(c fiber.Ctx) error {
	err := c.Next()
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(service.ErrNotFound, err):
		return c.Status(fiber.StatusNotFound).SendString(err.Error())

	case errors.Is(service.ErrAlreadyExists, err):
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	log.Err(err).Str("ip", c.IP()).Msg("error middleware")
	return c.SendStatus(fiber.StatusInternalServerError)
}

func authMiddleware(auth service.Auth) fiber.Handler {
	return func(c fiber.Ctx) error {
		token, ok := parseToken(c.Request())
		if !ok {
			log.Warn().Str("ip", c.IP()).Msg("auth middleware unauthorized access")
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		if auth.Validate(token) {
			return c.Next()
		}
		log.Warn().Str("ip", c.IP()).Msg("auth middleware invalid token")
		return c.SendStatus(fiber.StatusForbidden)
	}
}

func parseToken(r *fasthttp.Request) (string, bool) {
	header := string(r.Header.Peek(fiber.HeaderAuthorization))
	if header == "" {
		return "", false
	}
	token := strings.Split(header, "Bearer ")
	if len(token) != 2 {
		return "", false
	}
	return token[1], true
}
