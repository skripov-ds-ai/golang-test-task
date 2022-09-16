//go:build integration
// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"testing"
)

func TestGetAd(t *testing.T) {
	data := GetAdAPI{ID: 3, Fields: []string{}}
	bs, _ := json.Marshal(data)

	req, err := http.NewRequest("GET", "http://localhost:8888/get_ad", bytes.NewBuffer(bs))
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("resp Status:", resp.Status)
	body, _ := io.ReadAll(resp.Body)
	result := map[string]interface{}{}
	_ = json.Unmarshal(body, &result)
	fmt.Println("resp Body:", result)
}

func TestCreateAd(t *testing.T) {
	num, _ := decimal.NewFromString("123.76")
	b := make([]rune, 201)
	for i := range b {
		b[i] = 'a'
	}
	title := string(b)
	data := AdJSONItem{
		Title: title, Description: "xyz",
		Price:     num,
		ImageURLs: []string{},
	}
	bs, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", "http://localhost:8888/create_ad", bytes.NewBuffer(bs))
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("resp Status:", resp.Status)
	body, _ := io.ReadAll(resp.Body)
	result := map[string]interface{}{}
	_ = json.Unmarshal(body, &result)
	fmt.Println("resp Body:", result)
}
