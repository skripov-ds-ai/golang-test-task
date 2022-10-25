package entities

import "github.com/shopspring/decimal"

type (
	// CreateAdAnswer stores a status of created item
	// and an ID of created item(if it was created)
	CreateAdAnswer struct {
		Status ResultStatus `json:"status"`
		ID     *int         `json:"id"`
	}

	// AdJSONItem contains the information to write to AdItem table
	AdJSONItem struct {
		Title       string          `json:"title" validate:"required,min=1,max=200"`
		Description string          `json:"description" validate:"max=1000"`
		Price       decimal.Decimal `json:"price" validate:"required,numeric"`
		ImageURLs   []string        `json:"image_urls" validate:"max=3,checkURL"`
	}
)
