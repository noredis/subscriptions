package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func Logging(logger *zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		statusCode := c.Response().StatusCode()

		evt := logger.Info()
		if statusCode >= 500 {
			evt = logger.Error()
		} else if statusCode >= 400 {
			evt = logger.Warn()
		}

		evt.
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", statusCode).
			Dur("duration", duration).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Msg("request_completed")

		return err
	}
}
