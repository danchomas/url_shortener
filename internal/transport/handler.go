package transport

import (
	"url_shorter/internal/entity"
	"url_shorter/internal/service"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RateLimitMiddleware(c *fiber.Ctx) error {
	ip := c.IP()

	allowed, err := h.service.CheckRateLimit(ip)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	if !allowed {
		return c.Status(429).JSON(fiber.Map{"error": "Превышен лимит запросов"})
	}

	return c.Next()
}

func (h *Handler) Register(app *fiber.App) {
	app.Post("/shorten", h.RateLimitMiddleware, h.Shorten)

	app.Get("/r/:code", h.Redirect)
	app.Get("/stats/:code", h.Stats)
}

func (h *Handler) Shorten(c *fiber.Ctx) error {
	var req entity.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "невалидный JSON"})
	}

	if req.URL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "поле url обязательно"})
	}

	res, err := h.service.Shorten(req.URL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(res)
}

func (h *Handler) Redirect(c *fiber.Ctx) error {
	code := c.Params("code")
	url, err := h.service.GetOriginalURL(code)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "ссылка не найдена"})
	}
	return c.Redirect(url)
}

func (h *Handler) Stats(c *fiber.Ctx) error {
	code := c.Params("code")
	stats, err := h.service.GetStats(code)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "ссылка не найдена"})
	}
	return c.JSON(stats)
}
