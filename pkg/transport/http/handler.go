package http

import "github.com/gofiber/fiber/v3"

type Handler interface {
	Register(router fiber.Router)
}
