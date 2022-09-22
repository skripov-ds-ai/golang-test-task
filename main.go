package main

import (
	"fmt"
	"golang-test-task/database"
	"golang-test-task/facade"
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
	logic := facade.NewHandlerFacade(client, v, logger)

	mux := http.NewServeMux()
	endpoints := []string{"create_ad", "get_ad", "list_ads"}
	for _, endpoint := range endpoints {
		path := fmt.Sprintf("/api/v0.1/%s", endpoint)
		if h, ok := logic.GetHandler(endpoint); ok {
			mux.HandleFunc(path, h)
		} else {
			logger.Warn("handler endpoint does not contain in logic", zap.String("endpoint", endpoint))
		}
	}
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		panic(err)
	}
}
