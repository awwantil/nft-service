// service/kubo_service.go
package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ipfs/go-cid"
	"io"
	"main/internal/models"
	"mime/multipart"
	"net/http"
)

const kuboApiBaseUrl = "http://127.0.0.1:5001/api/v0"
const kuboGatewayUrlTemplate = "http://%s.ipfs.localhost:8080/"

// AddFileToIPFS загружает файл в узел Kubo и возвращает информацию о нем.
func AddFileToIPFS(fileHeader *multipart.FileHeader) (*models.AddResponse, string, string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, "", "", fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return nil, "", "", fmt.Errorf("не удалось создать form-file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, "", "", fmt.Errorf("не удалось скопировать данные файла: %w", err)
	}
	writer.Close()

	// Документация на Kubo RPC API подтверждает использование этого эндпоинта
	// Источник: https://docs.ipfs.tech/reference/kubo/rpc/
	req, err := http.NewRequest("POST", kuboApiBaseUrl+"/add", &requestBody)
	if err != nil {
		return nil, "", "", fmt.Errorf("не удалось создать запрос к Kubo: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", "", fmt.Errorf("ошибка при выполнении запроса к Kubo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, "", "", fmt.Errorf("Kubo API вернул ошибку: %s, тело ответа: %s", resp.Status, string(bodyBytes))
	}

	var addResp models.AddResponse
	if err := json.NewDecoder(resp.Body).Decode(&addResp); err != nil {
		return nil, "", "", fmt.Errorf("не удалось декодировать ответ от Kubo: %w", err)
	}

	// --- Начало новой логики ---
	// Декодируем полученный CIDv0 (начинается с "Qm")
	// Источник: https://pkg.go.dev/github.com/ipfs/go-cid#Decode
	cidV0, err := cid.Decode(addResp.Hash)
	if err != nil {
		return nil, "", "", fmt.Errorf("не удалось декодировать CID: %w", err)
	}

	cidV1 := cid.NewCidV1(cid.DagProtobuf, cidV0.Hash())
	gatewayURL := fmt.Sprintf(kuboGatewayUrlTemplate, cidV1.String())

	return &addResp, cidV1.String(), gatewayURL, nil
}

// PinCID закрепляет (pins) CID на узле Kubo.
func PinCID(cid string) (*models.PinResponse, error) {
	// Эндпоинт для закрепления: /api/v0/pin/add
	// Источник: https://github.com/ipfs/kubo
	url := fmt.Sprintf("%s/pin/add?arg=%s", kuboApiBaseUrl, cid)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос на закрепление: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса на закрепление: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Kubo API (pin) вернул ошибку: %s", resp.Status)
	}

	var pinResp models.PinResponse
	if err := json.NewDecoder(resp.Body).Decode(&pinResp); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ от Kubo (pin): %w", err)
	}
	return &pinResp, nil
}

// UnpinCID открепляет (unpins) CID с узла Kubo.
func UnpinCID(cid string) (*models.PinResponse, error) {
	// Эндпоинт для открепления: /api/v0/pin/rm
	// Источник: https://github.com/ipfs/kubo
	url := fmt.Sprintf("%s/pin/rm?arg=%s", kuboApiBaseUrl, cid)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос на открепление: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса на открепление: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Kubo API (unpin) вернул ошибку: %s", resp.Status)
	}

	var unpinResp models.PinResponse
	if err := json.NewDecoder(resp.Body).Decode(&unpinResp); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ от Kubo (unpin): %w", err)
	}
	return &unpinResp, nil
}

// ListPinnedCIDs возвращает список всех закрепленных CID.
func ListPinnedCIDs() (*models.PinLsResponse, error) {
	// Эндпоинт для получения списка закрепленных объектов: /api/v0/pin/ls
	// Источник: https://github.com/ipfs/kubo
	url := fmt.Sprintf("%s/pin/ls", kuboApiBaseUrl)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос на получение списка: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса на получение списка: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Kubo API (ls) вернул ошибку: %s", resp.Status)
	}

	var lsResp models.PinLsResponse
	if err := json.NewDecoder(resp.Body).Decode(&lsResp); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ от Kubo (ls): %w", err)
	}
	return &lsResp, nil
}
