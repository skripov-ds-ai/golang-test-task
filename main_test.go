//go:build integration
// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func Test_Main(t *testing.T) {
	data := GetAdAPI{ID: 0, Fields: []string{}}
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
