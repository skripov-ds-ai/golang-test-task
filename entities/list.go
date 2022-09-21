package entities

import "github.com/shopspring/decimal"

const PaginationSize int = 10

// Pagination contains all information required to paginate records' table
type Pagination struct {
	Offset int    `json:"offset"`
	By     string `json:"by"`
	Asc    bool   `json:"asc"`
}

// AdAPIListItem stores information is needed for pagination show
type AdAPIListItem struct {
	ID    int
	Title string
	Price decimal.Decimal
}

// ListAdsAnswer combine a status of a list process and a result
type ListAdsAnswer struct {
	Status string          `json:"status"`
	Result []AdAPIListItem `json:"result"`
}
