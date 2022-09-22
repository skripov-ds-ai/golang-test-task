package database

import (
	"github.com/shopspring/decimal"
	"reflect"
	"testing"
)

func TestCreateMapFromAdItemEmptyItem(t *testing.T) {
	item := AdItem{}
	var url *string
	expected := map[string]interface{}{"id": 0, "main_image_url": url, "title": "", "price": decimal.NewFromInt(0)}
	actual := item.CreateMap([]string{})

	if len(expected) != len(actual) {
		t.Fatalf("len(expected) != len(actual) ; expected = %v ; actual = %v", expected, actual)
	}
	//if !reflect.DeepEqual(expected, actual) {
	//	t.Fatalf("expected != actual ; expected = %v ; actual = %v", expected, actual)
	//}
	for k, v := range expected {
		val, ok := actual[k]
		if !ok {
			t.Fatalf("k = %s not found in actual ; expected = %v ; actual = %v", k, expected, actual)
		}
		if k == "price" {
			r1, ok1 := val.(decimal.Decimal)
			r2, ok2 := v.(decimal.Decimal)
			if !ok1 {
				t.Fatalf("expected[price] does not convert to decimal.Decimal ; expected = %v", expected)
			}
			if !ok2 {
				t.Fatalf("actual[price] does not convert to decimal.Decimal ; actual = %v", actual)
			}
			if !r1.Equal(r2) {
				t.Fatalf("expected[%s] != actual[%s] ; %v != %v ; expected = %v ; actual = %v", k, k, v, val, expected, actual)
			}
			continue
		}
		if v == nil && val != nil || v != nil && !reflect.DeepEqual(v, val) {
			t.Fatalf("expected[%s] != actual[%s] ; %v != %v ; expected = %v ; actual = %v", k, k, v, val, expected, actual)
		}
	}
}

func TestCreateMapFromAdItemNotNilMainImage(t *testing.T) {
	s := ""
	var url = &s
	mainURL := ImageURL{URL: *url}
	item := AdItem{MainImageURL: &mainURL}
	expected := map[string]interface{}{"id": 0, "main_image_url": url, "title": "", "price": decimal.NewFromInt(0)}
	actual := item.CreateMap([]string{})

	if len(expected) != len(actual) {
		t.Fatalf("len(expected) != len(actual) ; expected = %v ; actual = %v", expected, actual)
	}
	for k, v := range expected {
		val, ok := actual[k]
		if !ok {
			t.Fatalf("k = %s not found in actual ; expected = %v ; actual = %v", k, expected, actual)
		}
		if k == "price" {
			r1, ok1 := val.(decimal.Decimal)
			r2, ok2 := v.(decimal.Decimal)
			if !ok1 {
				t.Fatalf("expected[price] does not convert to decimal.Decimal ; expected = %v", expected)
			}
			if !ok2 {
				t.Fatalf("actual[price] does not convert to decimal.Decimal ; actual = %v", actual)
			}
			if !r1.Equal(r2) {
				t.Fatalf("expected[%s] != actual[%s] ; %v != %v ; expected = %v ; actual = %v", k, k, v, val, expected, actual)
			}
			continue
		}
		if r2, ok := val.(*string); ok {
			r1, ok := v.(*string)
			if ok {
				if *r1 == *r2 {
					continue
				}
				t.Fatalf("expected[%s] != actual[%s] ; %v != %v ; expected = %v ; actual = %v", k, k, *r1, *r2, expected, actual)
			}
		}
		if v == nil && val != nil || v != nil && v != val {
			t.Fatalf("expected[%s] != actual[%s] ; %v != %v ; expected = %v ; actual = %v", k, k, v, val, expected, actual)
		}
	}
}

func TestCreateMapFromAdItemEmptyItemAdditionalFieldsEmptyImageList(t *testing.T) {
	var url *string
	item := AdItem{}
	expected := map[string]interface{}{"id": 0, "main_image_url": url, "title": "", "description": "", "image_urls": []string{}, "price": decimal.NewFromInt(0)}
	actual := item.CreateMap([]string{"description", "image_urls"})

	if len(expected) != len(actual) {
		t.Fatalf("len(expected) != len(actual) ; expected = %v ; actual = %v", expected, actual)
	}
	for k, v := range expected {
		val, ok := actual[k]
		if !ok {
			t.Fatalf("k = %s not found in actual ; expected = %v ; actual = %v", k, expected, actual)
		}
		if k == "price" {
			r1, ok1 := val.(decimal.Decimal)
			r2, ok2 := v.(decimal.Decimal)
			if !ok1 {
				t.Fatalf("expected[price] does not convert to decimal.Decimal ; expected = %v", expected)
			}
			if !ok2 {
				t.Fatalf("actual[price] does not convert to decimal.Decimal ; actual = %v", actual)
			}
			if !r1.Equal(r2) {
				t.Fatalf("expected[%s] != actual[%s] ; %v != %v ; expected = %v ; actual = %v", k, k, v, val, expected, actual)
			}
			continue
		}
		if k == "image_urls" {
			r1, ok1 := v.([]string)
			r2, ok2 := val.([]string)
			if !ok1 || !ok2 || len(r1) != len(r2) {
				t.Fatalf("expected[%s] != actual[%s] ; %v != %v ; expected = %v ; actual = %v", k, k, v, val, expected, actual)
			}
		} else if v == nil && val != nil || v != nil && v != val {
			t.Fatalf("expected[%s] != actual[%s] ; %v != %v ; expected = %v ; actual = %v", k, k, v, val, expected, actual)
		}
	}
}
