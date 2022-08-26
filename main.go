package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	//"github.com/ericlagergren/decimal"
	"github.com/shopspring/decimal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

const paginationSize int = 10

type UniversalHandler struct {
	DB        *gorm.DB
	validator *validator.Validate
}

// ListAds is a function to get list of ads
//
// sorting by price/date_created; asc/desc order
// TODO: add pagination size to params
func (u *UniversalHandler) ListAds(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

// GetAd is a function to get concreate ad
//
// required fields: title, price, main_photo_url
// additional: by parameter `fields`(description, photo_urls)
func (u *UniversalHandler) GetAd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

// CreateAd is a function to create ad
//
// Params: title, description, photo_urls, price
// Return: ID of new ad, code of a result
func (u *UniversalHandler) CreateAd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

type ImageURL struct {
	gorm.Model
	URL  string
	AdID int
}

type AdJSONItem struct {
	Title       string          `json:"title" validate:"required"`
	Description string          `json:"description" validate:"required"`
	Price       decimal.Decimal `json:"price" validate:"required,numeric"`
	ImageURLs   []string        `json:"imageURLs" validate:"required"`
}

type AdItem struct {
	gorm.Model
	ID           int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Title        string
	Description  string
	Price        decimal.Decimal `sql:"type:decimal(20,8);"`
	ImageURLSs   []ImageURL      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	MainImageURL ImageURL        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (u *UniversalHandler) createAd(title, description string,
	photoURLs []string, price decimal.Decimal) (id int, err error) {
	size := 0
	if len(photoURLs) > 1 {
		size = len(photoURLs) - 1
	}

	var imageURLs = make([]ImageURL, size)

	if size > 0 {
		for i := range photoURLs {
			imageURLs[i].URL = photoURLs[i+1]
		}
	}

	mainImageURL := ImageURL{URL: photoURLs[0]}

	ad := AdItem{Title: title, Description: description, Price: price,
		ImageURLSs: imageURLs, MainImageURL: mainImageURL}

	db := u.DB.Create(&ad)
	err = db.Error
	if err != nil {
		return 0, err
	}
	return ad.ID, nil
}

func main() {
	// https://github.com/shopspring/decimal/issues/21
	decimal.MarshalJSONWithoutQuotes = true

	v := validator.New()

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	_ = UniversalHandler{DB: db, validator: v}

}
