package facade

import (
	"encoding/json"
	"golang-test-task/database"
	"golang-test-task/entities"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	// HandlerFacade is helper struct for getting handlers and hiding inner logic of an app
	HandlerFacade struct {
		dbClient  *database.Client
		validator *validator.Validate
		logger    *zap.Logger
		handlers  map[string]handler
	}
	handler   func(w http.ResponseWriter, r *http.Request)
	handlerBs func(w http.ResponseWriter, bs []byte)
)

// NewHandlerFacade is constructor for HandlerFacade
func NewHandlerFacade(dbClient *database.Client, validator *validator.Validate, logger *zap.Logger) *HandlerFacade {
	facade := HandlerFacade{dbClient: dbClient, validator: validator, logger: logger}
	facade.handlers = make(map[string]handler)
	facade.handlers["create_ad"] = facade.checkMethod("POST", facade.readAllWrap(facade.createAd))
	facade.handlers["get_ad"] = facade.checkMethod("GET", facade.readAllWrap(facade.getAd))
	facade.handlers["list_ads"] = facade.checkMethod("GET", facade.readAllWrap(facade.listAds))
	return &facade
}

// GetHandler gives a handler based on endpoint
func (hf *HandlerFacade) GetHandler(s string) (func(w http.ResponseWriter, r *http.Request), bool) {
	h, ok := hf.handlers[s]
	return h, ok
}

func (hf *HandlerFacade) readAllWrap(h handlerBs) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			hf.logger.Error("error during ReadAll")
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

func (hf *HandlerFacade) checkMethod(method string, h handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			hf.logger.Error("wrong method")
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

// listAds is a function to get list of ads
//
// sorting by price/date_created; asc/desc order
// TODO: add pagination size to params
// @Accept  json
// @Description An endpoint to get item list for pagination
// @Failure 400 {object} entities.ListAdsAnswer{}
// @Failure 500 {object} entities.ListAdsAnswer{}
// @Success 200 {object} entities.ListAdsAnswer{}
// @Router /list_ads [get]
func (hf *HandlerFacade) listAds(w http.ResponseWriter, bs []byte) {
	result := entities.ListAdsAnswer{Status: "error"}
	var pag entities.Pagination
	err := json.Unmarshal(bs, &pag)
	if err != nil {
		hf.logger.Error("error during Unmarshal in listAds", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	err = hf.validator.Struct(pag)
	if err != nil {
		hf.logger.Error("error during validating in listAds", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	// TODO: make paginationSize customizable
	items, err := hf.dbClient.ListAds(pag.Offset, entities.PaginationSize, pag.By, pag.Asc)
	if err != nil {
		hf.logger.Error("error during dbClient.listAds in listAds", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	itms := make([]map[string]interface{}, len(items))
	for i, v := range items {
		itms[i] = v.CreateMap()
	}
	result.Result = itms
	bs, _ = json.Marshal(result)
	_, _ = w.Write(bs)
}

// getAd is a function to get concreate ad
//
// required fields: title, price, main_photo_url
// additional: by parameter `fields`(description, photo_urls)
// @Accept  json
// @Description An endpoint to get item by id
// @Failure 400 {object} entities.GetAdAnswer{}
// @Failure 500 {object} entities.GetAdAnswer{}
// @Success 200 {object} entities.GetAdAnswer{}
// @Router /get_ad [get]
func (hf *HandlerFacade) getAd(w http.ResponseWriter, bs []byte) {
	result := entities.GetAdAnswer{Status: "error"}
	var api entities.GetAdAPI
	err := json.Unmarshal(bs, &api)
	if err != nil {
		hf.logger.Error("error during Unmarshal in getAd", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	err = hf.validator.Struct(api)
	if err != nil {
		hf.logger.Error("error during validating in getAd", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	item, err := hf.dbClient.GetAd(api.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		hf.logger.Error("error during dbClient.getAd in getAd", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}
	result.Status = "success"
	if item != nil {
		m := item.CreateMap(api.Fields)
		result.Result = &m
	}

	bs, _ = json.Marshal(result)
	_, _ = w.Write(bs)
}

// createAd is a function to create ad
//
// Params: title, description, photo_urls, price
// Return: ID of new ad, code of a result
// @Accept  json
// @Description An endpoint to create item
// @Failure 400 {object} entities.CreateAdAnswer{}
// @Failure 500 {object} entities.CreateAdAnswer{}
// @Success 200 {object} entities.CreateAdAnswer{}
// @Router /create_ad [post]
func (hf *HandlerFacade) createAd(w http.ResponseWriter, bs []byte) {
	result := entities.CreateAdAnswer{ID: nil, Status: "error"}
	var item entities.AdJSONItem
	err := json.Unmarshal(bs, &item)
	if err != nil {
		hf.logger.Error("error during Unmarshal in createAd", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}
	err = hf.validator.Struct(item)
	if err != nil {
		hf.logger.Error("error during validating in createAd", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		bs, _ = json.Marshal(result)
		_, _ = w.Write(bs)
		return
	}

	id, err := hf.dbClient.CreateAd(item)
	if err != nil {
		hf.logger.Error("error during dbClient.createAd in createAd", zap.Error(err))
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
