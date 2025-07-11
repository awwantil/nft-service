package models

// AddResponse представляет ответ от /api/v0/add
type AddResponse struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

// PinResponse представляет ответ от /api/v0/pin/add и /api/v0/pin/rm
type PinResponse struct {
	Pins []string `json:"Pins"`
}

// PinLsKey представляет информацию о закрепленном объекте
type PinLsKey struct {
	Type string `json:"Type"`
}

// PinLsResponse представляет ответ от /api/v0/pin/ls
type PinLsResponse struct {
	Keys map[string]PinLsKey `json:"Keys"`
}
