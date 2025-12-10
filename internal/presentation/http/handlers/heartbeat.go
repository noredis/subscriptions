package handlers

import "github.com/gofiber/fiber/v2"

type HeartbeatHandler struct{}

func NewHeartbeatHandler() *HeartbeatHandler {
	return &HeartbeatHandler{}
}

func (handler *HeartbeatHandler) Register(app *fiber.App) {
	app.Get("/heartbeat", handler.Heartbeat)
}

// Heartbeat проверяет доступность сервиса.
//
// @Summary      Проверка доступности
// @Description  Возвращает 200 OK, если сервис работает.
// @Tags         health
// @Success      200  "Сервис доступен"
// @Router       /heartbeat [get]
func (handler *HeartbeatHandler) Heartbeat(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
