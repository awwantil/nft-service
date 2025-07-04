// handler/file_handlers.go
package handler

import (
	"github.com/awwantil/nft-service/internal/service"
	
)

// UploadFileHandler обрабатывает загрузку файла.
func UploadFileHandler(c *fiber.Ctx) error {
	// Получаем файл из multipart-формы "file"
	// Источник: https://docs.gofiber.io/recipes/upload-file/
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Не удалось получить файл из формы",
			"data":    err.Error(),
		})
	}

	// Вызываем сервис для добавления файла в IPFS
	addResponse, err := service.AddFileToIPFS(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Ошибка при добавлении файла в IPFS",
			"data":    err.Error(),
		})
	}

	// Возвращаем успешный ответ с данными о файле
	// Источник: https://medium.com/@m7adeel/golang-backend-image-upload-api-887e07e5a70b
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Файл успешно загружен в IPFS",
		"data":    addResponse,
	})
}

// PinCidHandler обрабатывает закрепление CID.
func PinCidHandler(c *fiber.Ctx) error {
	cid := c.Params("cid")
	if cid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "CID не указан"})
	}

	pinResponse, err := service.PinCID(cid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(pinResponse)
}

// UnpinCidHandler обрабатывает открепление CID.
func UnpinCidHandler(c *fiber.Ctx) error {
	cid := c.Params("cid")
	if cid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "CID не указан"})
	}

	unpinResponse, err := service.UnpinCID(cid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(unpinResponse)
}

// ListPinsHandler обрабатывает запрос на получение списка закрепленных CID.
func ListPinsHandler(c *fiber.Ctx) error {
	lsResponse, err := service.ListPinnedCIDs()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(lsResponse)
}
