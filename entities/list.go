package entities

// PaginationSize is a constant for item list
const (
	PaginationSize int    = 10
	ByCreatedAt    string = "created_at"
	ByPrice        string = "price"
)

// ListAdsAnswer combine a status of a list process and a result
type ListAdsAnswer struct {
	Status string                   `json:"status"`
	Result []map[string]interface{} `json:"result"`
}
