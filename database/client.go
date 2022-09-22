package database

import (
	"fmt"

	"golang-test-task/entities"

	"gorm.io/gorm"
)

// Client is wrapper for database
type Client struct {
	db *gorm.DB
}

// NewClient creates custom db client
func NewClient(db *gorm.DB) *Client {
	c := Client{}
	c.db = db
	return &c
}

// GetAd gives an item and an error
func (c *Client) GetAd(id int) (res *AdItem, err error) {
	// TODO: add fields to use in .Select(fields)
	var item AdItem
	db := c.db.Preload("ImageURLs").Preload("MainImageURL").First(&item, id)
	err = db.Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// ListAds gives list of items
func (c *Client) ListAds(offset, paginationSize int, by string, asc bool) (resItems []entities.AdAPIListItem, err error) {
	var items []*entities.AdAPIListItem
	ascOrDesc := "asc"
	if !asc {
		ascOrDesc = "desc"
	}
	order := fmt.Sprintf("%s %s", by, ascOrDesc)
	db := c.db.Preload("ImageURLs").Preload("MainImageURL").Model(&AdItem{}).Limit(paginationSize).Offset(offset).Order(order).Find(&items)
	err = db.Error
	if err != nil {
		return resItems, err
	}
	for _, v := range items {
		resItems = append(resItems, *v)
	}
	return resItems, nil
}

// CreateAd creates an item
func (c *Client) CreateAd(item entities.AdJSONItem) (id int, err error) {
	var mainImageURL *ImageURL
	size := 0
	imgURLsSize := len(item.ImageURLs)
	if imgURLsSize > 0 {
		mainImageURL = &ImageURL{URL: item.ImageURLs[0]}
		if imgURLsSize > 1 {
			size = imgURLsSize - 1
		}
	}

	var imageURLs = make([]ImageURL, size)
	if imgURLsSize > 1 {
		arr := item.ImageURLs[1:]
		for i := range arr {
			if i == 0 {
				continue
			}
			imageURLs[i].URL = arr[i]
		}
	}

	ad := AdItem{Title: item.Title, Description: item.Description, Price: item.Price,
		ImageURLs: imageURLs, MainImageURL: mainImageURL,
	}

	db := c.db.Create(&ad)
	err = db.Error
	if err != nil {
		return 0, err
	}
	return ad.ID, nil
}
