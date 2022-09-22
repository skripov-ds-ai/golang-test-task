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

// @Title Ad Submitting Service
// @Version 0.1
// @Description A sample server for creating and getting ads

// @Contact.name Denis Skripov
// @Contact.email nizhikebinesi@gmail.com

// @Host localhost:8888
// @BasePath /api/v0.1

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
	mux.HandleFunc("/api/v0.1/create_ad", logic.checkMethod("POST", logic.readAllWrap(logic.CreateAd)))
	mux.HandleFunc("/api/v0.1/get_ad", logic.checkMethod("GET", logic.readAllWrap(logic.GetAd)))
	mux.HandleFunc("/api/v0.1/list_ads", logic.checkMethod("GET", logic.readAllWrap(logic.ListAds)))

	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		panic(err)
	}
}
