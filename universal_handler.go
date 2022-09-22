package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang-test-task/database"
	"golang-test-task/entities"
	"gorm.io/gorm"
	"io"
	"net/http"
)

type (
	universalHandler struct {
		dbClient  *database.Client
		validator *validator.Validate
		logger    *zap.Logger
	}
	handler   func(w http.ResponseWriter, r *http.Request)
	handlerBs func(w http.ResponseWriter, bs []byte)
)

func (u *universalHandler) readAllWrap(h handlerBs) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			u.logger.Error("error during ReadAll")
			result := make(map[string]interface{})
			result["status"] = "error"
			bs, _ = json.Marshal(result)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(bs)
			return
		}
		h(w, bs)
	}
}

func (u *universalHandler) checkMethod(method string, h handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			u.logger.Error("wrong method")
			result := make(map[string]interface{})
			result["status"] = "error"
			bs, _ := json.Marshal(result)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(bs)
			return
		}
		h(w, r)
	}
}

// ListAds is a function to get list of ads
//
// sorting by price/date_created; asc/desc order
// TODO: add pagination size to params
// @Accept  json
// @Description An endpoint to get item list for pagination
// @Failure 400 {object} entities.ListAdsAnswer{}
// @Failure 500 {object} entities.ListAdsAnswer{}
// @Success 200 {object} entities.ListAdsAnswer{}
// @Router /list_ads [get]
func (u *universalHandler) ListAds(w http.ResponseWriter, bs []byte) {
	result := entities.ListAdsAnswer{Status: "error"}
	var pag entities.Pagination
	err := json.Unmarshal(bs, &pag)
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
// @Accept  json
// @Description An endpoint to get item by id
// @Failure 400 {object} entities.GetAdAnswer{}
// @Failure 500 {object} entities.GetAdAnswer{}
// @Success 200 {object} entities.GetAdAnswer{}
// @Router /get_ad [get]
func (u *universalHandler) GetAd(w http.ResponseWriter, bs []byte) {
	result := entities.GetAdAnswer{Status: "error"}
	var api entities.GetAdAPI
	err := json.Unmarshal(bs, &api)
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
// @Accept  json
// @Description An endpoint to create item
// @Failure 400 {object} entities.CreateAdAnswer{}
// @Failure 500 {object} entities.CreateAdAnswer{}
// @Success 200 {object} entities.CreateAdAnswer{}
// @Router /create_ad [post]
func (u *universalHandler) CreateAd(w http.ResponseWriter, bs []byte) {
	result := entities.CreateAdAnswer{ID: nil, Status: "error"}
	var item entities.AdJSONItem
	err := json.Unmarshal(bs, &item)
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
