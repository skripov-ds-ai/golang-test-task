package entities

import "github.com/shopspring/decimal"

// PaginationSize is a constant for item list
const (
	PaginationSize int    = 10
	ByCreatedAt    string = "created_at"
	ByPrice        string = "price"
)

// ListAdsAnswer combine a status of a list process and a result
type (
	ListAdsAnswer struct {
		Status string          `json:"status"`
		Result []APIAdListItem `json:"result"`
	}
	APIAdListItem struct {
		ID           int             `json:"id"`
		Title        string          `json:"title"`
		Price        decimal.Decimal `json:"price"`
		MainImageURL *string         `json:"main_image_url"`
	}
)
