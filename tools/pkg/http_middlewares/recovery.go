package httpmiddlewares

import (
	"runtime"

	"github.com/gofiber/fiber/v2"

	"main/tools/pkg/logger"
)

// NewRecoveryMiddleware panic recover middleware
func NewRecoveryMiddleware(logger *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				logger.Error("recovering from err", "error", err, "stack", buf)

				c.SendStatus(fiber.StatusInternalServerError)
				c.Write([]byte(`{"error":"internal error occured"}`))
			}
		}()
		return c.Next()
	}
}
