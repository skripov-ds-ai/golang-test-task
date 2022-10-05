package entities

import "github.com/shopspring/decimal"

const (
	// PaginationSize is a constant for item list
	PaginationSize int = 10
	// ByCreatedAt is a constant for sorting by a column of item creation
	ByCreatedAt string = "created_at"
	// ByPrice is a constant for sorting by a column of item price
	ByPrice string = "price"
)

type (
	// ListAdsAnswer combine a status of a list process and a result
	ListAdsAnswer struct {
		Status string          `json:"status"`
		Result []APIAdListItem `json:"result"`
	}

	// APIAdListItem is a struct to bind data for listAd
	APIAdListItem struct {
		ID           int             `json:"id"`
		Title        string          `json:"title"`
		Price        decimal.Decimal `json:"price"`
		MainImageURL *string         `json:"main_image_url"`
	}
)
