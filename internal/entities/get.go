package entities

// GetAdAnswer couples a status of processing and a request's result
type GetAdAnswer struct {
	Status string                  `json:"status"`
	Result *map[string]interface{} `json:"result"`
}
