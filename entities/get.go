package entities

// GetAdAPI stores information about interesting item and which fields to show
type GetAdAPI struct {
	ID     int      `json:"id" validate:"numeric,min=0"`
	Fields []string `json:"fields"`
}

// GetAdAnswer couples a status of processing and a request's result
type GetAdAnswer struct {
	Status string                  `json:"status"`
	Result *map[string]interface{} `json:"result"`
}
