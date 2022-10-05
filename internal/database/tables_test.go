package database

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestAdAPIListItem_CreateMap(t *testing.T) {
	item := AdListItem{}
	var url *string
	zero := decimal.NewFromInt(0)
	expected := map[string]interface{}{"id": 0, "main_image_url": url, "title": "", "price": zero}
	actual := item.CreateMap()

	if len(expected) != len(actual) {
		t.Fatalf("len(expected) != len(actual) ; expected = %v ; actual = %v", expected, actual)
	}
	price, ok := actual["price"]
	if !ok {
		t.Fatalf("price key does not contain in expected ; expected = %v", expected)
	}
	actualPrice, ok := price.(decimal.Decimal)
	if !ok {
		t.Fatalf("price from expected does not convert to decimal.Decimal ; expected = %v", expected)
	}
	if !actualPrice.Equal(zero) {
		t.Fatalf("prices of expected and actual is not equal ; expected = %v ; actual = %v", expected, actualPrice)
	}
	delete(expected, "price")
	delete(actual, "price")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected != actual ; expected = %v ; actual = %v", expected, actual)
	}
}

func TestAdItem_CreateMapEmptyItem(t *testing.T) {
	item := AdItem{}
	var url *string
	zero := decimal.NewFromInt(0)
	expected := map[string]interface{}{"id": 0, "main_image_url": url, "title": "", "price": decimal.NewFromInt(0)}
	actual := item.CreateMap([]string{})

	if len(expected) != len(actual) {
		t.Fatalf("len(expected) != len(actual) ; expected = %v ; actual = %v", expected, actual)
	}
	price, ok := actual["price"]
	if !ok {
		t.Fatalf("price key does not contain in expected ; expected = %v", expected)
	}
	actualPrice, ok := price.(decimal.Decimal)
	if !ok {
		t.Fatalf("price from expected does not convert to decimal.Decimal ; expected = %v", expected)
	}
	if !actualPrice.Equal(zero) {
		t.Fatalf("prices of expected and actual is not equal ; expected = %v ; actual = %v", expected, actualPrice)
	}
	delete(expected, "price")
	delete(actual, "price")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected != actual ; expected = %v ; actual = %v", expected, actual)
	}
}

func TestAdItem_CreateMapWithMainImage(t *testing.T) {
	s := ""
	var url = &s
	zero := decimal.NewFromInt(0)
	mainURL := ImageURL{URL: *url}
	item := AdItem{MainImageURL: &mainURL}
	expected := map[string]interface{}{"id": 0, "main_image_url": url, "title": "", "price": decimal.NewFromInt(0)}
	actual := item.CreateMap([]string{})

	if len(expected) != len(actual) {
		t.Fatalf("len(expected) != len(actual) ; expected = %v ; actual = %v", expected, actual)
	}
	if len(expected) != len(actual) {
		t.Fatalf("len(expected) != len(actual) ; expected = %v ; actual = %v", expected, actual)
	}
	price, ok := actual["price"]
	if !ok {
		t.Fatalf("price key does not contain in expected ; expected = %v", expected)
	}
	actualPrice, ok := price.(decimal.Decimal)
	if !ok {
		t.Fatalf("price from expected does not convert to decimal.Decimal ; expected = %v", expected)
	}
	if !actualPrice.Equal(zero) {
		t.Fatalf("prices of expected and actual is not equal ; expected = %v ; actual = %v", expected, actualPrice)
	}
	delete(expected, "price")
	delete(actual, "price")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected != actual ; expected = %v ; actual = %v", expected, actual)
	}
}

func TestAdItem_CreateMapEmptyItemAdditionalFieldsEmptyImageList(t *testing.T) {
	var url *string
	zero := decimal.NewFromInt(0)
	item := AdItem{}
	expected := map[string]interface{}{"id": 0, "main_image_url": url, "title": "", "description": "", "image_urls": []string{}, "price": decimal.NewFromInt(0)}
	actual := item.CreateMap([]string{"description", "image_urls"})

	if len(expected) != len(actual) {
		t.Fatalf("len(expected) != len(actual) ; expected = %v ; actual = %v", expected, actual)
	}
	if len(expected) != len(actual) {
		t.Fatalf("len(expected) != len(actual) ; expected = %v ; actual = %v", expected, actual)
	}
	price, ok := actual["price"]
	if !ok {
		t.Fatalf("price key does not contain in expected ; expected = %v", expected)
	}
	actualPrice, ok := price.(decimal.Decimal)
	if !ok {
		t.Fatalf("price from expected does not convert to decimal.Decimal ; expected = %v", expected)
	}
	if !actualPrice.Equal(zero) {
		t.Fatalf("prices of expected and actual is not equal ; expected = %v ; actual = %v", expected, actualPrice)
	}
	delete(expected, "price")
	delete(actual, "price")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected != actual ; expected = %v ; actual = %v", expected, actual)
	}
}
