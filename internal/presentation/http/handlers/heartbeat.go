package handlers

import "github.com/gofiber/fiber/v2"

type HeartbeatHandler struct{}

func NewHeartbeatHandler() *HeartbeatHandler {
	return &HeartbeatHandler{}
}

func (handler *HeartbeatHandler) Register(app *fiber.App) {
	app.Get("/heartbeat", handler.Heartbeat)
}

func (handler *HeartbeatHandler) Heartbeat(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
