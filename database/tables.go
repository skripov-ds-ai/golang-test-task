package database

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ImageURL is a table for image URLs for items
type ImageURL struct {
	gorm.Model
	URL      string
	AdItemID int
}

// AdItem is a table which contains information about all items
type AdItem struct {
	gorm.Model
	ID           int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Title        string
	Description  string
	Price        decimal.Decimal `sql:"type:decimal(20,8);"`
	ImageURLs    []ImageURL      `gorm:"foreignKey:AdItemID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	MainImageURL *ImageURL       `gorm:"foreignKey:AdItemID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// CreateMap creates map based on fields to show
func (item *AdItem) CreateMap(fields []string) (m map[string]interface{}) {
	m = map[string]interface{}{}
	m["id"] = item.ID
	m["title"] = item.Title
	var u *string
	if item.MainImageURL != nil {
		u = &item.MainImageURL.URL
	}
	m["main_image_url"] = u
	for _, field := range fields {
		if field == "description" {
			m["description"] = item.Description
		} else if field == "image_urls" {
			imgUrls := make([]string, 0)
			for _, v := range item.ImageURLs {
				imgUrls = append(imgUrls, v.URL)
			}
			m["image_urls"] = imgUrls
		}
	}
	return m
}

// AdAPIListItem stores information is needed for pagination show
type AdAPIListItem struct {
	ID           int
	Title        string
	Price        decimal.Decimal
	MainImageURL *ImageURL `gorm:"foreignKey:AdItemID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// CreateMap creates map based on fields to show
func (item *AdAPIListItem) CreateMap() (m map[string]interface{}) {
	m = map[string]interface{}{}
	m["id"] = item.ID
	m["title"] = item.Title
	var u *string
	if item.MainImageURL != nil {
		u = &item.MainImageURL.URL
	}
	m["main_image_url"] = u
	return m
}
