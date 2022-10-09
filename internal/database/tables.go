package database

import (
	"golang-test-task/internal/entities"

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
func (item *AdItem) CreateMap(fields []string) entities.APIAdItem {
	itm := entities.APIAdItem{}
	itm.ID = item.ID
	itm.Title = item.Title
	itm.Price = item.Price
	var u *string
	if item.MainImageURL != nil {
		u = &item.MainImageURL.URL
	}
	itm.MainImageURL = u
	for _, field := range fields {
		if field == "description" {
			itm.Description = item.Description
		} else if field == "image_urls" {
			imgUrls := make([]string, 0)
			for _, v := range item.ImageURLs {
				imgUrls = append(imgUrls, v.URL)
			}
			itm.ImageURLs = imgUrls
		}
	}
	return itm
}

// CreateMapFromFields creates map based on fields' map to show
func (item *AdItem) CreateMapFromFields(fields map[string]struct{}) entities.APIAdItem {
	itm := entities.APIAdItem{}
	itm.ID = item.ID
	itm.Title = item.Title
	itm.Price = item.Price
	var u *string
	if item.MainImageURL != nil {
		u = &item.MainImageURL.URL
	}
	itm.MainImageURL = u
	if _, ok := fields["description"]; ok {
		itm.Description = item.Description
	}
	if _, ok := fields["image_urls"]; ok {
		imgUrls := make([]string, 0)
		for _, v := range item.ImageURLs {
			imgUrls = append(imgUrls, v.URL)
		}
		itm.ImageURLs = imgUrls
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
	itm.ID = item.ID
	itm.Title = item.Title
	itm.Price = item.Price
	var u *string
	if item.MainImageURL != nil {
		u = &item.MainImageURL.URL
	}
	itm.MainImageURL = u
	return itm
}
