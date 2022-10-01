package entities

// PaginationSize is a constant for item list
const PaginationSize int = 10

// Pagination contains all information required to paginate records' table
type Pagination struct {
	Offset int    `json:"offset" validate:"required,min=0"`
	By     string `json:"by" validate:"required,oneof=price created_at"`
	Asc    bool   `json:"asc"`
}

// ListAdsAnswer combine a status of a list process and a result
type ListAdsAnswer struct {
	Status string                   `json:"status"`
	Result []map[string]interface{} `json:"result"`
}
