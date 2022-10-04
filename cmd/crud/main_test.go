package main

// func (a *App) executeRequest(req *http.Request) *httptest.ResponseRecorder {
//	res := httptest.NewRecorder()
//	a.Router.ServeHTTP(res, req)
//	return res
// }

// import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"github.com/shopspring/decimal"
//	"io"
//	"net/http"
//	"testing"
//
//	"golang-test-task/entities"
// )
//
// func TestGetAd(t *testing.T) {
//	id := 100
//	data := entities.GetAdAPI{ID: id, Fields: []string{"description", "image_urls"}}
//	bs, _ := json.Marshal(data)
//
//	req, err := http.NewRequest("GET", "http://localhost:8888/api/v0.1/get_ad", bytes.NewBuffer(bs))
//	if err != nil {
//		panic(err)
//	}
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		panic(err)
//	}
//	defer resp.Body.Close()
//
//	fmt.Println("resp Status:", resp.Status)
//	body, _ := io.ReadAll(resp.Body)
//	result := map[string]interface{}{}
//	_ = json.Unmarshal(body, &result)
//	fmt.Println("resp Body:", result)
// }
//
// func TestListAd(t *testing.T) {
//	//id := 100
//	data := entities.Pagination{By: "id"}
//	bs, _ := json.Marshal(data)
//
//	req, err := http.NewRequest("POST", "http://localhost:8888/api/v0.1/list_ads", bytes.NewBuffer(bs))
//	if err != nil {
//		panic(err)
//	}
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		panic(err)
//	}
//	defer resp.Body.Close()
//
//	fmt.Println("resp Status:", resp.Status)
//	body, _ := io.ReadAll(resp.Body)
//	result := map[string]interface{}{}
//	_ = json.Unmarshal(body, &result)
//	fmt.Println("resp Body:", result)
// }
//
// func TestCreateAd(t *testing.T) {
//	num, _ := decimal.NewFromString("123.76")
//	b := make([]rune, 200)
//	for i := range b {
//		b[i] = 'e'
//	}
//	title := string(b)
//	data := entities.AdJSONItem{
//		Title: title, Description: "xyz",
//		Price:     num,
//		ImageURLs: []string{"http://yandex.ru"},
//	}
//	bs, _ := json.Marshal(data)
//
//	req, err := http.NewRequest("POST", "http://localhost:8888/api/v0.1/create_ad", bytes.NewBuffer(bs))
//	if err != nil {
//		panic(err)
//	}
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		panic(err)
//	}
//	defer resp.Body.Close()
//
//	fmt.Println("resp Status:", resp.Status)
//	body, _ := io.ReadAll(resp.Body)
//	result := map[string]interface{}{}
//	_ = json.Unmarshal(body, &result)
//	fmt.Println("resp Body:", result)
// }
