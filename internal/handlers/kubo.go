package handlers

import (
	"github.com/gofiber/fiber/v2"
	"main/internal/service"
	"main/tools/pkg/logger"
)

// KuboHandlers
type KuboHandlers struct {
	logger *logger.Logger
}

// NewAuthHandlers конструктор для обработчиков IDM методов
func NewKuboHandlers(logger *logger.Logger) *KuboHandlers {
	return &KuboHandlers{
		logger: logger,
	}
}

// UploadFileHandler обрабатывает загрузку файла.
func UploadFileHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Не удалось получить файл из формы",
			"data":    err.Error(),
		})
	}

	// Вызываем обновленный сервис, который возвращает больше данных
	addResponse, cidV1, gatewayURL, err := service.AddFileToIPFS(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Ошибка при добавлении файла в IPFS",
			"data":    err.Error(),
		})
	}

	// Формируем расширенный JSON-ответ.
	// Источник: https://dev.to/hackmamba/robust-media-upload-with-golang-and-cloudinary-fiber-version-2cmf
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Файл успешно загружен в IPFS",
		"data": fiber.Map{
			"name":       addResponse.Name,
			"size":       addResponse.Size,
			"cidV0":      addResponse.Hash,
			"cidV1":      cidV1,
			"gatewayUrl": gatewayURL,
		},
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
