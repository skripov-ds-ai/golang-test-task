package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"

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

type Pagination struct {
	Offset int    `json:"offset"`
	By     string `json:"by"`
	Asc    bool   `json:"asc"`
}

// ListAds is a function to get list of ads
//
// sorting by price/date_created; asc/desc order
// TODO: add pagination size to params
func (u *UniversalHandler) ListAds(w http.ResponseWriter, r *http.Request) {
	result := ListAdsAnswer{Status: "error"}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	var pag Pagination
	err = json.Unmarshal(bs, &pag)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	// TODO: make paginationSize customizable
	items, err := u.listAds(pag.Offset, paginationSize, pag.By, pag.Asc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	result.Result = items
	bs1, _ := json.Marshal(result)
	_, _ = w.Write(bs1)
}

func (u *UniversalHandler) listAds(offset, paginationSize int, by string, asc bool) (resItems []AdAPIListItem, err error) {
	items := []*AdAPIListItem{}
	ascOrDesc := "asc"
	if !asc {
		ascOrDesc = "desc"
	}
	order := fmt.Sprintf("%s %s", by, ascOrDesc)
	db := u.DB.Model(&AdItem{}).Limit(paginationSize).Offset(offset).Order(order).Find(&items)
	err = db.Error
	if err != nil {
		return resItems, err
	}
	for _, v := range items {
		resItems = append(resItems, *v)
	}
	return resItems, nil
}

// GetAd is a function to get concreate ad
//
// required fields: title, price, main_photo_url
// additional: by parameter `fields`(description, photo_urls)
func (u *UniversalHandler) GetAd(w http.ResponseWriter, r *http.Request) {
	result := GetAdAnswer{Status: "error"}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	var api GetAdAPI
	err = json.Unmarshal(bs, &api)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	item, err := u.getAd(api.ID, api.Fields)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}
	result.Status = "success"
	if item != nil {
		m := createMapFromAdItem(*item, api.Fields)
		result.Result = &m
	}

	// TODO: add item handling
	fmt.Println(item)
	bs1, _ := json.Marshal(result)
	_, _ = w.Write(bs1)
}

func createMapFromAdItem(item AdItem, fields []string) (m map[string]interface{}) {
	m = map[string]interface{}{}
	m["id"] = item.ID
	m["title"] = item.Title
	m["main_image_url"] = item.MainImageURL.URL
	for _, field := range fields {
		if field == "description" {
			m["description"] = item.Description
		} else if field == "" {
			imgUrls := make([]string, 0)
			for _, v := range item.ImageURLSs {
				imgUrls = append(imgUrls, v.URL)
			}
			m["image_urls"] = imgUrls
		}
	}
	return m
}

type GetAdAPI struct {
	ID     int      `json:"id"`
	Fields []string `json:"fields"`
}

type GetAdAnswer struct {
	Status string                  `json:"status"`
	Result *map[string]interface{} `json:"result"`
}

func (u *UniversalHandler) getAd(id int, fields []string) (res *AdItem, err error) {
	var item AdItem
	db := u.DB.First(&item, id)
	err = db.Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// CreateAd is a function to create ad
//
// Params: title, description, photo_urls, price
// Return: ID of new ad, code of a result
func (u *UniversalHandler) CreateAd(w http.ResponseWriter, r *http.Request) {
	result := CreateAdAnswer{ID: nil, Status: "error"}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	var item AdJSONItem
	err = json.Unmarshal(bs, &item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}
	err = u.validator.Struct(item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}
	//for _, e := range err.(validator.ValidationErrors) {
	//	fmt.Println(e)
	//}

	id, err := u.createAd(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}
	result.Status = "success"
	result.ID = &id

	bs1, _ := json.Marshal(result)
	_, _ = w.Write(bs1)
}

type CreateAdAnswer struct {
	Status string `json:"status"`
	ID     *int   `json:"id"`
}

type ListAdsAnswer struct {
	Status string          `json:"status"`
	Result []AdAPIListItem `json:"result"`
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

type AdAPIListItem struct {
	ID    int
	Title string
	Price decimal.Decimal
}

func (u *UniversalHandler) createAd(item AdJSONItem) (id int, err error) {
	size := 0
	if len(item.ImageURLs) > 1 {
		size = len(item.ImageURLs) - 1
	}

	var imageURLs = make([]ImageURL, size)

	if size > 0 {
		for i := range item.ImageURLs {
			imageURLs[i].URL = item.ImageURLs[i+1]
		}
	}

	mainImageURL := ImageURL{URL: item.ImageURLs[0]}

	ad := AdItem{Title: item.Title, Description: item.Description, Price: item.Price,
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
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	logic := UniversalHandler{DB: db, validator: v}

	mux := http.NewServeMux()
	mux.HandleFunc("/create_ad", logic.CreateAd)
	mux.HandleFunc("/get_ad", logic.GetAd)
	mux.HandleFunc("/list_ads", logic.ListAds)

	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		panic(err)
	}
}
