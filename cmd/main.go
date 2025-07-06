// main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"main/internal/handlers"
)

func main() {
	app := fiber.New()

	api := app.Group("/api/v1")

	api.Post("/files", handler.UploadFileHandler)

	// Маршруты для управления закреплением (pin)
	api.Post("/pins/:cid", handler.PinCidHandler)
	api.Delete("/pins/:cid", handler.UnpinCidHandler)
	api.Get("/pins", handler.ListPinsHandler)

	// Источник: https://medium.com/@m7adeel/golang-backend-image-upload-api-887e07e5a70b
	log.Fatal(app.Listen(":3000"))
}
