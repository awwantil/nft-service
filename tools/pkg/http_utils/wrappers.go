package httputils

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"

	tvoerrors "main/tools/pkg/tvo_errors"
)

// FiberJSONWrapper wrapper for response json
func FiberJSONWrapper(callback func(c *fiber.Ctx) (interface{}, error)) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		res, err := callback(c)
		if err != nil {
			return HandleError(c, FiberStatusByErr(err), err)
		}

		return json.NewEncoder(c.Type("json", "utf-8").Response().BodyWriter()).Encode(res)
	}
}

// FiberStatusByErr returns the appropriate HTTP status code for the given error.
func FiberStatusByErr(err error) int {
	switch {
	case errors.Is(err, tvoerrors.ErrInvalidRequestData):
		return fiber.StatusBadRequest
	case errors.Is(err, tvoerrors.ErrNotFound):
		return fiber.StatusNotFound
	case errors.Is(err, tvoerrors.ErrUnauthorized):
		return fiber.StatusUnauthorized
	case errors.Is(err, tvoerrors.ErrForbidden):
		return fiber.StatusForbidden
	case errors.Is(err, tvoerrors.ErrConflict):
		return fiber.StatusConflict
	default:
		return fiber.StatusInternalServerError
	}
}
