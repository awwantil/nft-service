// main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/your-username/go-kubo-service/handler"
)

func main() {
	// Создаем новый экземпляр Fiber
	// Источник: https://medium.com/@m7adeel/golang-backend-image-upload-api-887e07e5a70b
	app := fiber.New()

	// Группа маршрутов API
	api := app.Group("/api/v1")

	// Маршрут для загрузки файла в IPFS
	// Источник: https://withcodeexample.com/file-upload-handling-golang-fiber-guide/
	api.Post("/files", handler.UploadFileHandler)

	// Маршруты для управления закреплением (pin)
	// Источник: https://github.com/ipfs/kubo - документация на RPC API
	api.Post("/pins/:cid", handler.PinCidHandler)
	api.Delete("/pins/:cid", handler.UnpinCidHandler)
	api.Get("/pins", handler.ListPinsHandler)

	// Запуск сервера на порту 3000
	// Источник: https://medium.com/@m7adeel/golang-backend-image-upload-api-887e07e5a70b
	log.Fatal(app.Listen(":3000"))
}
