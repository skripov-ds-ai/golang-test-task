package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const paginationSize int = 10

type universalHandler struct {
	DB        *gorm.DB
	validator *validator.Validate
	logger    *zap.Logger
}

// Pagination contains all information required to paginate records' table
type Pagination struct {
	Offset int    `json:"offset"`
	By     string `json:"by"`
	Asc    bool   `json:"asc"`
}

// ListAds is a function to get list of ads
//
// sorting by price/date_created; asc/desc order
// TODO: add pagination size to params
func (u *universalHandler) ListAds(w http.ResponseWriter, r *http.Request) {
	result := ListAdsAnswer{Status: "error"}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		u.logger.Error("error during ReadAll in ListAds", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	var pag Pagination
	err = json.Unmarshal(bs, &pag)
	if err != nil {
		u.logger.Error("error during Unmarshal in ListAds", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	// TODO: make paginationSize customizable
	items, err := u.listAds(pag.Offset, paginationSize, pag.By, pag.Asc)
	if err != nil {
		u.logger.Error("error during listAds in ListAds", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	result.Result = items
	bs1, _ := json.Marshal(result)
	_, _ = w.Write(bs1)
}

func (u *universalHandler) listAds(offset, paginationSize int, by string, asc bool) (resItems []AdAPIListItem, err error) {
	items := []*AdAPIListItem{}
	ascOrDesc := "asc"
	if !asc {
		ascOrDesc = "desc"
	}
	order := fmt.Sprintf("%s %s", by, ascOrDesc)
	db := u.DB.Preload("ImageURLs").Preload("MainImageURL").Model(&AdItem{}).Limit(paginationSize).Offset(offset).Order(order).Find(&items)
	err = db.Error
	if err != nil {
		u.logger.Error("error during getting data in listAd", zap.Error(err))
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
func (u *universalHandler) GetAd(w http.ResponseWriter, r *http.Request) {
	result := GetAdAnswer{Status: "error"}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		u.logger.Error("error during ReadAll in GetAd", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	var api GetAdAPI
	err = json.Unmarshal(bs, &api)
	if err != nil {
		u.logger.Error("error during Unmarshal in GetAd", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	item, err := u.getAd(api.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.logger.Error("error during getAd in GetAd", zap.Error(err))
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
	// fmt.Println(item)
	bs1, _ := json.Marshal(result)
	_, _ = w.Write(bs1)
}

func createMapFromAdItem(item AdItem, fields []string) (m map[string]interface{}) {
	m = map[string]interface{}{}
	m["id"] = item.ID
	m["title"] = item.Title
	var url *string = nil
	if item.MainImageURL != nil {
		// fmt.Println(item)
		url = &item.MainImageURL.URL
	}
	// fmt.Println(item)
	m["main_image_url"] = url
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

// GetAdAPI stores information about interesting item and which fields to show
type GetAdAPI struct {
	ID     int      `json:"id"`
	Fields []string `json:"fields"`
}

// GetAdAnswer couples a status of processing and a request's result
type GetAdAnswer struct {
	Status string                  `json:"status"`
	Result *map[string]interface{} `json:"result"`
}

func (u *universalHandler) getAd(id int) (res *AdItem, err error) {
	// TODO: add fields to use in .Select(fields)
	var item AdItem
	db := u.DB.Preload("ImageURLs").Preload("MainImageURL").First(&item, id).Association("MainImageURL")
	err = db.Error
	if err != nil {
		u.logger.Error("error during getting data in getAd", zap.Error(err))
		return nil, err
	}
	return &item, nil
}

// CreateAd is a function to create ad
//
// Params: title, description, photo_urls, price
// Return: ID of new ad, code of a result
func (u *universalHandler) CreateAd(w http.ResponseWriter, r *http.Request) {
	result := CreateAdAnswer{ID: nil, Status: "error"}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		u.logger.Error("error during ReadAll in CreateAd", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}

	var item AdJSONItem
	err = json.Unmarshal(bs, &item)
	if err != nil {
		u.logger.Error("error during Unmarshal in CreateAd", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}
	err = u.validator.Struct(item)
	if err != nil {
		u.logger.Error("error during validating in CreateAd", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		bs1, _ := json.Marshal(result)
		_, _ = w.Write(bs1)
		return
	}
	// for _, e := range err.(validator.ValidationErrors) {
	//	fmt.Println(e)
	//}

	id, err := u.createAd(item)
	if err != nil {
		u.logger.Error("error during createAd in CreateAd", zap.Error(err))
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

// CreateAdAnswer stores a status of created item
// and an ID of created item(if it was created)
type CreateAdAnswer struct {
	Status string `json:"status"`
	ID     *int   `json:"id"`
}

// ListAdsAnswer combine a status of a list process and a result
type ListAdsAnswer struct {
	Status string          `json:"status"`
	Result []AdAPIListItem `json:"result"`
}

// ImageURL is a table for image URLs for items
type ImageURL struct {
	gorm.Model
	URL      string
	AdItemID int
}

// AdJSONItem contains the information to write to AdItem table
type AdJSONItem struct {
	Title       string          `json:"title" validate:"required,min=1,max=200"`
	Description string          `json:"description" validate:"required,max=1000"`
	Price       decimal.Decimal `json:"price" validate:"required,numeric"`
	ImageURLs   []string        `json:"imageURLs" validate:"required,max=3,checkURL"`
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

// AdAPIListItem stores information is needed for pagination show
type AdAPIListItem struct {
	ID    int
	Title string
	Price decimal.Decimal
}

func (u *universalHandler) createAd(item AdJSONItem) (id int, err error) {
	var mainImageURL *ImageURL = nil
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

	fmt.Println("!", mainImageURL)

	ad := AdItem{Title: item.Title, Description: item.Description, Price: item.Price,
		ImageURLs: imageURLs, MainImageURL: mainImageURL,
	}

	db := u.DB.Create(&ad)
	err = db.Error
	if err != nil {
		u.logger.Error("error during creating datum in createAd", zap.Error(err))
		return 0, err
	}
	return ad.ID, nil
}

func main() {
	// https://github.com/shopspring/decimal/issues/21
	decimal.MarshalJSONWithoutQuotes = true

	v := validator.New()
	_ = v.RegisterValidation("checkURL", func(fl validator.FieldLevel) bool {
		arr, ok := fl.Field().Interface().([]string)
		if !ok {
			return false
		}
		for _, a := range arr {
			_, err := url.ParseRequestURI(a)
			if err != nil {
				return false
			}
		}
		return true
	})

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// c := zap.NewDevelopmentConfig()
	// c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// l, _ := c.Build()
	// defer l.Sync()
	// loggerG := zapgorm2.New(l)
	// loggerG.SetAsDefault() // optional: configure gorm to use this zapgorm.Logger for callbacks

	dsn := "host=postgres user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: loggerG})
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Logger: glogger.Default.LogMode(glogger.Warn),
	})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&AdItem{}, &ImageURL{})
	if err != nil {
		panic("failed to automigrate")
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()

	// logger, _ := zap.NewDevelopment()
	defer func() {
		_ = logger.Sync()
	}()

	// d, err := db.DB()
	////d.Close()
	// d.Ping()
	// if err != nil {
	//	fmt.Println(err)
	//}

	logic := universalHandler{DB: db, validator: v, logger: logger}

	mux := http.NewServeMux()
	mux.HandleFunc("/create_ad", logic.CreateAd)
	mux.HandleFunc("/get_ad", logic.GetAd)
	mux.HandleFunc("/list_ads", logic.ListAds)

	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		panic(err)
	}
}
