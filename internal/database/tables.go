package database

import (
	"github.com/shopspring/decimal"
	"golang-test-task/internal/entities"
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
func (item *AdItem) CreateMap(fields []string) entities.APIAdItem {
	itm := entities.APIAdItem{}
	//m = map[string]interface{}{}
	itm.ID = item.ID
	//m["id"] = item.ID
	itm.Title = item.Title
	//m["title"] = item.Title
	itm.Price = item.Price
	//m["price"] = item.Price
	var u *string
	if item.MainImageURL != nil {
		u = &item.MainImageURL.URL
	}
	itm.MainImageURL = u
	//m["main_image_url"] = u
	for _, field := range fields {
		if field == "description" {
			itm.Description = item.Description
			//m["description"] = item.Description
		} else if field == "image_urls" {
			imgUrls := make([]string, 0)
			for _, v := range item.ImageURLs {
				imgUrls = append(imgUrls, v.URL)
			}
			itm.ImageURLs = imgUrls
			//m["image_urls"] = imgUrls
		}
	}
	return itm
}

// AdListItem stores information is needed for pagination show
type AdListItem struct {
	ID           int
	Title        string
	Price        decimal.Decimal
	MainImageURL *ImageURL `gorm:"foreignKey:AdItemID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// CreateMap creates map based on fields to show
func (item *AdListItem) CreateMap() entities.APIAdListItem {
	itm := entities.APIAdListItem{}
	//m = map[string]interface{}{}
	itm.ID = item.ID
	//m["id"] = item.ID
	itm.Title = item.Title
	//m["title"] = item.Title
	itm.Price = item.Price
	//m["price"] = item.Price
	var u *string
	if item.MainImageURL != nil {
		u = &item.MainImageURL.URL
	}
	itm.MainImageURL = u
	//m["main_image_url"] = u
	return itm
}
