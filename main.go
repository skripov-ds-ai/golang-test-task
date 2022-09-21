package main

import (
	"encoding/json"
	"golang-test-task/database"
	"golang-test-task/entities"
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

type universalHandler struct {
	dbClient  *database.Client
	validator *validator.Validate
	logger    *zap.Logger
}

// ListAds is a function to get list of ads
//
// sorting by price/date_created; asc/desc order
// TODO: add pagination size to params
func (u *universalHandler) ListAds(w http.ResponseWriter, r *http.Request) {
	result := entities.ListAdsAnswer{Status: "error"}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		bs, _ := json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		u.logger.Error("error during ReadAll in ListAds", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	var pag entities.Pagination
	err = json.Unmarshal(bs, &pag)
	if err != nil {
		u.logger.Error("error during Unmarshal in ListAds", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	// TODO: make paginationSize customizable
	items, err := u.dbClient.ListAds(pag.Offset, entities.PaginationSize, pag.By, pag.Asc)
	if err != nil {
		u.logger.Error("error during dbClient.ListAds in ListAds", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	result.Result = items
	bs, _ = json.Marshal(result)
	_, _ = w.Write(bs)
}

// GetAd is a function to get concreate ad
//
// required fields: title, price, main_photo_url
// additional: by parameter `fields`(description, photo_urls)
func (u *universalHandler) GetAd(w http.ResponseWriter, r *http.Request) {
	result := entities.GetAdAnswer{Status: "error"}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		bs, _ := json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		u.logger.Error("error during ReadAll in GetAd", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	var api entities.GetAdAPI
	err = json.Unmarshal(bs, &api)
	if err != nil {
		u.logger.Error("error during Unmarshal in GetAd", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	item, err := u.dbClient.GetAd(api.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		u.logger.Error("error during dbClient.GetAd in GetAd", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}
	result.Status = "success"
	if item != nil {
		m := item.CreateMapFromAdItem(api.Fields)
		result.Result = &m
	}

	bs, _ = json.Marshal(result)
	_, _ = w.Write(bs)
}

// CreateAd is a function to create ad
//
// Params: title, description, photo_urls, price
// Return: ID of new ad, code of a result
func (u *universalHandler) CreateAd(w http.ResponseWriter, r *http.Request) {
	result := entities.CreateAdAnswer{ID: nil, Status: "error"}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		bs, _ := json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		u.logger.Error("error during ReadAll in CreateAd", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	var item entities.AdJSONItem
	err = json.Unmarshal(bs, &item)
	if err != nil {
		u.logger.Error("error during Unmarshal in CreateAd", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}
	err = u.validator.Struct(item)
	if err != nil {
		u.logger.Error("error during validating in CreateAd", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	id, err := u.dbClient.CreateAd(item)
	if err != nil {
		u.logger.Error("error during dbClient.CreateAd in CreateAd", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}
	result.Status = "success"
	result.ID = &id

	bs, _ = json.Marshal(result)
	_, _ = w.Write(bs)
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

	dsn := "host=postgres user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&database.AdItem{}, &database.ImageURL{})
	if err != nil {
		panic("failed to automigrate")
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	defer func() {
		_ = logger.Sync()
	}()

	client := database.NewClient(db)
	logic := universalHandler{dbClient: client, validator: v, logger: logger}

	mux := http.NewServeMux()
	mux.HandleFunc("/create_ad", logic.CreateAd)
	mux.HandleFunc("/get_ad", logic.GetAd)
	mux.HandleFunc("/list_ads", logic.ListAds)

	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		panic(err)
	}
}
